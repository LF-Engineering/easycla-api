# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Script to perform database migration from dynamodb instance in easyclav1 to RDS in easyclav2
"""

import logging
import os
import sys
import threading
from datetime import datetime
from queue import Queue
from threading import Thread
from typing import Dict, List


import click
import psycopg2

import log
from cla import utils
from utils import PostgresSchema
from cla.models.dynamo_models import (
    Company,
    Gerrit,
    GitHubOrg,
    Project,
    Repository,
    Signature,
    User,
)

IMPORT_THREADS = 4
successful_import = 0
failed_import = 0
tables_queues = Queue()
models_map = {
    User: "users",
    Company: "companies",
    Repository: "repositories",
    GitHubOrg: "github_orgs",
    Gerrit: "gerrit_instances",
    Project: "projects",
    Signature: "signatures",
}
logger = log.setup_custom_logger("dynamodb", prefix="data-migration-{}".format(os.environ.get("STAGE")),)

models = [User, Company, Repository, GitHubOrg, Gerrit, Project, Signature]


success = 0
failed = 0




def get_table(model):
    """ 
    Function that returns the associative table for the given instance 

    :param model: A dynamodb model instance
    """
    if type(model) in models:
        return models_map[type(model)]
    else:
        raise Exception("Failed to get the corresponding Dynamodb Table for the instance")


def update_queue(models):
    """
    Utility function that adds models to be migrated to the queue

    :param models: Models are enqueued for migration process in the worker threads
    """
    for item in models:
        tables_queues.put(item)


def import_to_pg(connection, cursor, queue, dry_run):
    """
    Worker that processes insertion based on the tables queue

    :param connection: Database connector
    :param cursor: database cursor instance
    :queue: The queue that stores models to be migrated
    :dry_run: flag that will save models if set to True
    """
    while True:
        item = tables_queues.get()
        if item is None:
            break
        # Parse cla.models and update associated table
        try:
            table = get_table(item)
            schema = PostgresSchema(item)
            attributes = schema.get_tables(table)
            # attributes = table_attributes(table,item)

            # table, insert_dict = get_table_attributes(item)

            cols = ""
            values = []
            for key, value in attributes().items():
                if value is not None:
                    cols += ",{}".format(key)
                    values.append(value)

            placeholders = ["%s" for _ in range(len(values))]
            str_placeholders = ", ".join(placeholders)

            insert_sql = "INSERT INTO {}({}) VALUES ({})".format(table, cols[1:], str_placeholders)

            if not dry_run:
                cursor.execute(insert_sql, tuple(values))
                connection.commit()
                logger.info("Run insert sql: {},{}".format(insert_sql,tuple(values)))
            else:
                logger.info("{},{} ".format(insert_sql, tuple(values)))
            


        except (Exception, psycopg2.Error) as error:
            logger.error(error)
            connection.rollback()
            
        queue.task_done()


@click.command()
@click.option(
    "--aws-region", is_flag=False, type=click.STRING, default="us-east-1", help="the aws region",
)
@click.option("--rds-database", is_flag=False, type=click.STRING, help="the rds database")
@click.option(
    "--rds-host", is_flag=False, type=click.STRING, default="localhost", help="the rds host",
)
@click.option(
    "--rds-username", is_flag=False, type=click.STRING, default="postgres", help="the rds username",
)
@click.option(
    "--rds-password", is_flag=False, type=click.STRING, default="postgres", help="the rds password",
)
@click.option(
    "--rds-port", is_flag=False, type=click.STRING, default="5432", help="the database port",
)
@click.option(
    "--dry-run",
    is_flag=True,
    default=False,
    help="Dry run of the migration process. Data is not saved in the postgres db",
)
def main(
    aws_region, rds_database, rds_host, rds_username, rds_password, rds_port, dry_run,
):
    """
    This script runs data migration from dynamodb to postgresql
    """

    if os.environ.get("STAGE") is None:
        logger.warning("Please set the 'STAGE' environment varaible - typically one of: {dev, staging, prod}")
        return

    stage = os.environ.get("STAGE")

    if dry_run:
        exit_cmd = input(
            "This is a dry run. "
            "You are running the script for the '{}' environment. "
            'Press <ENTER> to continue ("exit" to exit): '.format(stage)
        )
    else:
        exit_cmd = input(
            "This is NOT a dry run. "
            "You are running the script for the '{}' environment. "
            'Press <ENTER> to continue ("exit" to exit): '.format(stage)
        )
    if exit_cmd == "exit":
        return

    start_time = datetime.now()
    lock = threading.Lock()
    threads = []
    count = 0
    for _ in range(IMPORT_THREADS):
        connection = psycopg2.connect(
            database=rds_database, host=rds_host, user=rds_username, password=rds_password, port=rds_port,
        )
        cursor = connection.cursor()
        worker = Thread(target=import_to_pg, args=(connection, cursor, tables_queues, dry_run))
        worker.start()
        threads.append(worker)

    # TODO: cla models for signature,company-invitations and user-permissions
    try:
        user = utils.get_user_instance()
        company = utils.get_company_instance()
        project = utils.get_project_instance()
        repository = utils.get_repository_instance()
        github_orgs = utils.get_github_organization_instance()
        gerrit = utils.get_gerrit_instance()
        #signature = utils.get_signature_instance()
        update_queue(user.all())
        update_queue(company.all())
        update_queue(project.all())
        #update_queue(signature.all())
        update_queue(repository.all())
        update_queue(gerrit.all())
        update_queue(github_orgs.all())
        # block until all tasks are done
        tables_queues.join()
    except Exception as err:
        logger.error(err)

    # End workers
    for _ in range(IMPORT_THREADS):
        tables_queues.put(None)
    for thread in threads:
        thread.join()

    duration = datetime.now() - start_time
    logger.info("Data migration to the pg database run for a duration of {} ".format(duration))


if __name__ == "__main__":
    sys.exit(main())
