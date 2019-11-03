# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Routine to set a specific SSM parameter value.
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
@click.option('--name', is_flag=False, type=click.STRING,
              help='the parameter key/name')
@click.option('--value', is_flag=False, type=click.STRING,
              help='the parameter value')
@click.option('--aws-region', is_flag=False, type=click.STRING,
              default='us-east-1',
              help='the AWS region - default is us-east-1')
@click.option('--log-dir', is_flag=False, type=click.STRING,
              default='.',
              help='the log output folder - default is the current folder')
@click.option('--overwrite', is_flag=True, help='over write values')
@click.option('-v', '--verbose', is_flag=True, help='verbose flag')
def main(name, value, aws_region, log_dir, overwrite, verbose):
    """
    Routine to set a specific SSM parameter value.
    """
    if not os.path.isdir(log_dir):
        print(f'Log directory does not exist: \'{log_dir}\' '
              '- please create or adjust --log-dir parameter')
        return

    log_level = INFO
    if verbose:
        log_level = DEBUG

    overwrite_flag = False
    if overwrite:
        overwrite_flag = True

    if not name:
        print(f'Missing parameter \'--name\' - please set on the command line')
        return
    if not value:
        print(f'Missing parameter \'--value\' - please set on the command line')
        return

    start_time = datetime.now()
    logger = log.setup_custom_logger(
        'root',
        log_dir=log_dir,
        level=log_level,
        prefix=f'set-parameter-{os.environ.get("STAGE")}')

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
    logger.info(f'log-dir         : {log_dir}')
    logger.info(f'log level       : {getLevelName(log_level)}')

    try:
        ssm_client.put_parameter(
            Name=name,
            Value=value,
            Type='String',
            Overwrite=overwrite_flag
        )
    except ClientError as e:
        logger.info(f'Skipping {name}: {value} - '
                    f'key already exists and overwrite set to: {overwrite_flag}')

    logger.info('Finished export - duration: {}'.format(datetime.now() - start_time))


if __name__ == "__main__":
    sys.exit(main())
