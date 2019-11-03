# Tool Notes

## Install Python

Install a version of python 3.7, if not already installed on your system.

One approach is to use a python version manager to all your system to support
multiple python versions. `pipenv` is one such tool.

To install the python version manager, run:

```bash
# Mac
brew install pipenv
```

Once installed, you can list the  Python versions that are available.

```bash
pyenv install --list | grep " 3\.[678]"
```

Then install a specific version of your choice:

```bash
# This will take some time
pyenv install 3.8.0
```

After which, you can list the versions on your system:

```bash
pyenv versions
```

Output might look like:

```code
  system
  3.6.9
* 3.7.4 (set by /Users/ddeal/.pyenv/version)
  3.8.0
```

The * indicates that the system Python version is active currently. You’ll also
notice that this is set by a file in your root pyenv directory. This means that,
by default, you are still using your system Python.

If, for example, you wanted to use version 3.8.0, then you can use the global
command:

```bash
# Set this as the version of python
pyenv global 3.8.0
python3 -V
Python 3.8.0
```

If you ever want to go back to the system version of Python as the default, you
can run this:

```bash
pyenv global system
python -V
Python 2.7.12

pyenv versions
```

## Create a Virtual Environment

Once python is installed, create a virtual environment/sandbox for this project.
You may need to install `virtualenv` by running `pip3 install virtualenv`.

```bash
which python3
virtualenv --python=python3 .venv
```

Then simply load the environment/sandbox:

```bash
source .venv/bin/activate

# Verify:
python3 --version
```

Once we are loaded, we can install the application/tool dependencies:

```bash
pip3 install -r requirements
```

## Setup AWS

```bash
# Set the staging environment, typically one of: dev, staging, prod
export STAGE=dev
# Make sure you have your AWS profile setup and export the profile you want to use
export AWS_PROFILE="insert your aws profile"

# setup the virtual env again, if not already done:
source .venv/bin/activate
pip3 install -r requirements.txt
```

## Running the Tools

Below are the usage and examples for running the tools. The `ssm-export.py` and
`ssm-import.py` are useful for exporting and importing SSM parameters across
AWS regions.

### SSM Export

```bash
python3 ssm-export.py -h
Usage: ssm-export.py [OPTIONS]

  Routine to export the AWS SSM parameters to a JSON document suitable for
  subsequent importing.

Options:
  --output-filename TEXT  the output filename for the export - default is
                          output.json
  --aws-region TEXT       the AWS region - default is us-east-1
  --log-dir TEXT          the log output folder - default is the current
                          folder
  -v, --verbose           verbose flag
  -h, --help              Show this message and exit.
```

Example:

```bash
python3 ssm-export.py -v --output-filename output.json --aws-region us-east-1
```

### SSM Import

```bash
python3 ssm-import.py -h
Usage: ssm-import.py [OPTIONS]

  Routine to import a list of key/value pairs from a JSON document to AWS
  SSM

Options:
  --input-filename TEXT  the input filename for the import - default is
                         input.json
  --aws-region TEXT      the AWS region - default is us-east-1
  --dry-run              flag to indicate if this is a dry run, when set would
                         not upload the parameters to SSM
  --log-dir TEXT         the log output folder - default is the current folder
  --overwrite            over write values
  -v, --verbose          verbose flag
  -h, --help             Show this message and exit.
```

Example:

```bash
python3 ssm-import.py --input-filename output.json --aws-region us-east-2 -v
```

### SSM Get Parameter

Usage:

```bash
python3 ssm-get-parameter.py -h                                               
Usage: ssm-get-parameter.py [OPTIONS]

  Routine to get a specific SSM parameter value.

Options:
  --name TEXT        the parameter key/name
  --aws-region TEXT  the AWS region - default is us-east-1
  --log-dir TEXT     the log output folder - default is the current folder
  -v, --verbose      verbose flag
  -h, --help         Show this message and exit.
```

Example:

```bash
python3 ssm-get-parameter.py -v --name cla-rds-host-dev --aws-region us-east-2
```

### SSM Set Parameter

Usage:

```bash
python3 ssm-set-parameter.py -h
Usage: ssm-set-parameter.py [OPTIONS]

  Routine to set a specific SSM parameter value.

Options:
  --name TEXT        the parameter key/name
  --value TEXT       the parameter value
  --aws-region TEXT  the AWS region - default is us-east-1
  --log-dir TEXT     the log output folder - default is the current folder
  --overwrite        over write values
  -v, --verbose      verbose flag
  -h, --help         Show this message and exit.
```

Example:

```bash
python3 ssm-set-parameter.py -v --name cla-rds-host-dev --value cla-rds-cluster-dev.cluster-abcabcabcab.us-east-2.rds.amazonaws.com --aws-region us-east-2 --overwrite
```

## Lint

```bash
pylint *.py
```

## Unit Tests

```bash
python3 -m unittest *-test.py
```
