import logging
import sys
from datetime import datetime
from queue import Queue
from threading import Thread
from typing import List

import psycopg2
import click

from cla import utils
from cla.models.dynamo_models import (
    Company,
    CompanyInviteModel,
    Gerrit,
    GitHubOrg,
    Project,
    Repository,
    Signature,
    User,
)

IMPORT_THREADS = 4
tables_queues = Queue()


def import_to_pg():
    #TODO: Parse the dynamodb tables and do the export to postgresql
    while True:
        item = tables_queues.get()
        if item is None:
            break
        #Process records in dynamo table and migrate to postgresql
        print('Port Table to postgres')
        tables_queues.task_done()
        print('Processed table')

@click.command()
@click.option(
    "--aws-region",
    is_flag=False,
    type=click.STRING,
    default="us-east-1",
    help="the aws region",
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
    "--rds-port", is_flag=False, type=click.STRING, default="5432", help="database port"
)
@click.option("--verbose", is_flag=True, help="verbose")
@click.option(
    "--dry-run",
    is_flag=True,
    help="Dry run of the migration process. Data is not saved in the postgres db",
)
def main(aws_region, rds_host, rds_username, rds_password, rds_port, verbose, dry_run):
    """
    This script runs data migration from dynamodb to postgresql
    """
    start_time = datetime.now()
    company = utils.get_company_instance()
    companies = company.all()
    print('Size of companies: {}'.format(len(companies)))
    threads = []
    for _ in range(IMPORT_THREADS):
        worker = Thread(target=import_to_pg )
        worker.start()
        threads.append(worker)

    # TODO: get all tables to be exported. 
    company = utils.get_company_instance()
    companies = company.all()
    for comp in companies:
        tables_queues.put(comp)

    # block until all tasks are done
    tables_queues.join()

    # End workers
    for _ in range(IMPORT_THREADS):
        tables_queues.put(None)
    for thread in threads:
        thread.join()
    
    duration = datetime.now() - start_time
    print('Data migration to the pg database run for a duration of {} '.format(duration))

if __name__ == "__main__":
    sys.exit(main())
