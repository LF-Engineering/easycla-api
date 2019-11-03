# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Routine to import a list of key/value pairs from a JSON document to AWS SSM
"""
from __future__ import annotations

import json
import os
import sys
from datetime import datetime
from logging import INFO, DEBUG, getLevelName

import boto3
from botocore.exceptions import ClientError

import click
import log


@click.command(context_settings={'help_option_names': ['-h', '--help']})
@click.option('--input-filename', is_flag=False, type=click.STRING,
              default='input.json',
              help='the input filename for the import - default is input.json')
@click.option('--aws-region', is_flag=False, type=click.STRING,
              default='us-east-1',
              help='the AWS region - default is us-east-1')
@click.option('--dry-run', is_flag=True,
              help=('flag to indicate if this is a dry run, when set would not'
                    'upload the parameters to SSM'))
@click.option('--log-dir', is_flag=False, type=click.STRING,
              default='.',
              help='the log output folder - default is the current folder')
@click.option('--overwrite', is_flag=True, help='over write values')
@click.option('-v', '--verbose', is_flag=True, help='verbose flag')
def main(input_filename, aws_region, dry_run, log_dir, overwrite, verbose):
    """
    Routine to import a list of key/value pairs from a JSON document to AWS SSM
    """
    if not os.path.isdir(log_dir):
        print(f'Log directory does not exist: \'{log_dir}\' '
              '- please create or adjust --log-dir parameter')
        return

    if not os.path.isfile(input_filename):
        print(f'Input filename does not exist: \'{input_filename}\' '
              '- please specify a valid input file.')
        return

    log_level = INFO
    if verbose:
        log_level = DEBUG

    overwrite_flag = False
    if overwrite:
        overwrite_flag = True

    start_time = datetime.now()
    logger = log.setup_custom_logger(
        'root',
        log_dir=log_dir,
        level=log_level,
        prefix=f'import-{os.environ.get("STAGE")}')

    if os.environ.get('STAGE') is None:
        logger.warning('Please set \'STAGE\' environment variable.')
        return
    stage = os.environ.get('STAGE')

    aws_session = boto3.session.Session(
        profile_name='lfproduct-{}'.format(stage),
        region_name=aws_region)
    ssm_client = aws_session.client('ssm')

    logger.info(f'STAGE           : {stage}')
    logger.info(f'AWS REGION      : {aws_region}')
    logger.info(f'dry-run         : {dry_run}')
    logger.info(f'log-dir         : {log_dir}')
    logger.info(f'log level       : {getLevelName(log_level)}')

    with open(input_filename) as json_file:
        data = json.load(json_file)

        for kv in data:
            logger.info(f'Processing {kv["Name"]}: {kv["Value"]}')
            if dry_run:
                logger.info('Skipping upload - dry-run mode')
                continue

            try:
                ssm_client.put_parameter(
                    Name=kv["Name"],
                    Value=kv["Value"],
                    Type='String',
                    Overwrite=overwrite_flag
                )
            except ClientError as e:
                logger.info(f'Skipping {kv["Name"]}: {kv["Value"]} - '
                            f'key already exists and overwrite set to: {overwrite_flag}')

    logger.info('Finished export - duration: {}'.format(datetime.now() - start_time))


if __name__ == "__main__":
    sys.exit(main())
