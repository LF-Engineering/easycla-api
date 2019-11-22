# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Script to perform database migration from dynamodb instance in easyclav1 to RDS in easyclav2
"""

import logging
import sys
from datetime import datetime
from queue import Queue
from threading import Thread
import os
from typing import List

import psycopg2
import click

from cla import utils
from cla.models.dynamo_models import Company, User

IMPORT_THREADS = 4
tables_queues = Queue()


def update_queue(models):
    """
    Utility function that adds models to the queue
    """
    for item in models:
        tables_queues.put(item)


def import_to_pg(connection, cursor, queue, dry_run):
    """
    Worker that processes insertion based on the tables queue
    """
    while True:
        item = tables_queues.get()
        if item is None:
            break
            # Parse cla.models and export to postgresql

        try:
            if isinstance(item, User):
                user_sql = "INSERT INTO users(user_id,lf_email,\
                              lf_username,\
                              user_github_id,\
                              user_company_id,\
                              user_name,\
                              user_github_name\
                              ) VALUES (%s,%s,%s,%s,%s,%s,%s)\
                           "
                cursor.execute(
                    user_sql,
                    (
                        item.model.user_id,
                        item.model.lf_email,
                        item.model.lf_username,
                        item.model.user_github_id,
                        item.model.user_company_id,
                        item.model.user_name,
                        item.model.user_github_username,
                    )
                )
                connection.commit()
                logging.info("{} saved ".format(item.model.user_name))
        except (Exception, psycopg2.Error) as db_error:
            logging.error(db_error)
            connection.rollback()

        queue.task_done()


@click.command()
@click.option(
    "--aws-region",
    is_flag=False,
    type=click.STRING,
    default="us-east-1",
    help="the aws region",
)
@click.option(
    "--rds-database", is_flag=False, type=click.STRING, help="the rds database"
)
@click.option(
    "--rds-host",
    is_flag=False,
    type=click.STRING,
    default="localhost",
    help="the rds host",
)
@click.option(
    "--rds-username",
    is_flag=False,
    type=click.STRING,
    default="postgres",
    help="the rds username",
)
@click.option(
    "--rds-password",
    is_flag=False,
    type=click.STRING,
    default="postgres",
    help="the rds password",
)
@click.option(
    "--rds-port",
    is_flag=False,
    type=click.STRING,
    default="5432",
    help="the database port",
)
@click.option("--verbose", is_flag=True, help="verbose")
@click.option(
    "--dry-run",
    is_flag=True,
    default=False,
    help="Dry run of the migration process. Data is not saved in the postgres db",
)
def main(
    aws_region,
    rds_database,
    rds_host,
    rds_username,
    rds_password,
    rds_port,
    verbose,
    dry_run,
):
    """
    This script runs data migration from dynamodb to postgresql
    """
    if os.environ.get("STAGE") is None:
        logging.warning(
            "Please set the 'STAGE' environment varaible - typically one of: {dev, staging, prod}"
        )
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
    threads = []
    for _ in range(IMPORT_THREADS):
        connection = psycopg2.connect(
            database=rds_database,
            host=rds_host,
            user=rds_username,
            password=rds_password,
            port=rds_port,
        )
        cursor = connection.cursor()
        worker = Thread(
            target=import_to_pg, args=(connection, cursor, tables_queues, dry_run)
        )
        worker.start()
        threads.append(worker)

    # TODO: get all tables to be exported.
    company = utils.get_company_instance()
    user = utils.get_user_instance()
    update_queue(user.all())
    update_queue(company.all())

    """
    project = utils.get_project_instance()
    user = utils.get_user_instance()
    user_permission = utils.UserPermissions()
    print("User Permissions!")
    #print(user_permission.all())
    repository = utils.get_repository_instance()
    github_org = utils.get_github_organization_instance()
    gerrit = utils.get_gerrit_instance()
    """

    # block until all tasks are done
    tables_queues.join()

    # End workers
    for _ in range(IMPORT_THREADS):
        tables_queues.put(None)
    for thread in threads:
        thread.join()

    duration = datetime.now() - start_time
    print(
        "Data migration to the pg database run for a duration of {} ".format(duration)
    )


if __name__ == "__main__":
    sys.exit(main())
