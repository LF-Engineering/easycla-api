# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Routine to export the AWS SSM parameters to a JSON document suitable for subsequent importing.
"""
from __future__ import annotations

import json
import os
import sys
from datetime import datetime
from logging import Logger, INFO, DEBUG, getLevelName
from typing import List

import boto3
import click
import log


def get_ssm_keys(logger: Logger, ssm_client, key_filter: List[str]) -> List[str]:
    """
    Routine to fetch the SSM keys for the specified filter.
    :param logger: the logger to use
    :param ssm_client: the AWS SSM client object
    :param key_filter: the key name filter, e.g. ['cla']
    :return: a list of SSM keys
    """
    # Basically: aws ssm describe-parameters --filters "Key=Name,Values=cla" --profile lfproduct-dev
    logger.debug(f'Querying for the list of keys using filter: {filter}')
    response = ssm_client.describe_parameters(
        Filters=[{'Key': 'Name', 'Values': key_filter}],
        MaxResults=50)
    keys = []
    params = response['Parameters']
    for param in params:
        keys.append(param['Name'])
        logger.debug(f'Key: {param["Name"]}')

    # If we have more (pagination)
    while 'NextToken' in response:
        logger.debug('Querying for next page of keys...')
        response = ssm_client.describe_parameters(
            NextToken=response['NextToken'],
            Filters=[{'Key': 'Name', 'Values': key_filter}],
            MaxResults=50)

        params = response['Parameters']
        for param in params:
            keys.append(param['Name'])
            logger.debug(f'Key: {param["Name"]}')

    return keys


def get_ssm_values_batch(logger: Logger, ssm_client, keys) -> List[dict]:
    """
    Fetches the SSM key/value pairs for the specified key list. This is
    typically done in batches of size 10 per the API restrictions.

    :param logger: the system logger
    :param ssm_client: the AWS SSM client object
    :param keys: a list of SSM keys to use as the fetch/query
    :return: a list of key/value pairs containing the SSM parameter values
    """

    output = []
    ssm_dict = {'WithDecryption': False, 'Names': keys}
    # response = ssm_client.get_parameters_by_path(Path='/cla', **ssm_dict)
    logger.debug(f'Querying for the list of keys/values using key list: {keys}')
    response = ssm_client.get_parameters(**ssm_dict)
    params = response['Parameters']
    for param in params:
        # Example:
        #  {'Name':'my-key',
        #   'Type':'String',
        #   'Value':'my-value',
        #   'Version':1,
        #   'LastModifiedDate':datetime.datetime(2019,1,9,9,37,47,304000,tzinfo=tzlocal()),
        #   'ARN':'arn:aws:ssm:us-east-1:XXXXXXXXXXXX:parameter/salesforce-security-token-lfx-dev'}
        output.append({'Name': param['Name'], 'Value': param['Value']})

    # Manual pagination, since boto doesn't support it yet for get_parameters_by_path
    while 'NextToken' in response:
        logger.debug('Querying for next page of keys/values...')
        response = ssm_client.get_parameters(NextToken=response['NextToken'], **ssm_dict)
        params = response['Parameters']
        for param in params:
            output.append({'Name': param['Name'], 'Value': param['Value']})

    # logger.debug(f'Returning result: {output}')
    return output


def get_ssm_values(logger: Logger, ssm_client, keys: List[str]) -> List[dict]:
    """
    Fetches the SSM key/value pairs for the specified key list. This is
    routine brakes down the list into sizable batches, invokes the helper
    to fetch the results and aggregates the results.

    :param logger: the system logger
    :param ssm_client: the AWS SSM client object
    :param keys: a list of SSM keys to use as the fetch/query
    :return: a list of key/value pairs containing the SSM parameter values
    """
    output = []

    while len(keys) > 0:
        key_batch = []
        for _ in range(0, min(10, len(keys))):
            key_batch.append(keys.pop(0))

        output.append(get_ssm_values_batch(logger, ssm_client, key_batch))

    # Flatten and return the list
    return [y for x in output for y in x]


@click.command(context_settings={'help_option_names': ['-h', '--help']})
@click.option('--output-filename', is_flag=False, type=click.STRING,
              default='output.json',
              help='the output filename for the export - default is output.json')
@click.option('--aws-region', is_flag=False, type=click.STRING,
              default='us-east-1',
              help='the AWS region - default is us-east-1')
@click.option('--log-dir', is_flag=False, type=click.STRING,
              default='.',
              help='the log output folder - default is the current folder')
@click.option('-v', '--verbose', is_flag=True, help='verbose flag')
def main(output_filename, aws_region, log_dir, verbose):
    """
    Routine to export the AWS SSM parameters to a JSON document suitable for
    subsequent importing.
    """
    if not os.path.isdir(log_dir):
        print(f'Log directory does not exist: \'{log_dir}\' '
              '- please create or adjust --log-dir parameter')
        return

    log_level = INFO
    if verbose:
        log_level = DEBUG

    start_time = datetime.now()
    logger = log.setup_custom_logger(
        'root',
        log_dir=log_dir,
        level=log_level,
        prefix=f'export-{os.environ.get("STAGE")}')

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

    logger.debug('Querying for SSM parameters...')
    # https://github.com/shibboleth66/gorgonzola/blob/20f984b8dd28a388ce8e769fe9185b9af022c1db/gorgonzola/aws_ssm_global_parameters.py
    # https://www.programcreek.com/python/example/97943/boto3.client
    # https://github.com/elonmusk408/ansible1/blob/5db7501ebdacf65a8cf076da35ed6c3011c4c58a/lib/ansible/plugins/lookup/aws_ssm.py

    # Query for the list of keys matching our filter
    keys: List[str] = get_ssm_keys(logger, ssm_client, ['cla'])
    # Query for the key/value pairs from the key list
    output = get_ssm_values(logger, ssm_client, keys)
    # Save to the output file
    with open(output_filename, 'w') as outfile:
        json.dump(output, outfile, indent=2)

    logger.info(f'Wrote results to: {output_filename}')
    logger.info('Finished export - duration: {}'.format(datetime.now() - start_time))


if __name__ == "__main__":
    sys.exit(main())
