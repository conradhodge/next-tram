# Next tram

![build](https://github.com/conradhodge/next-tram/workflows/build/badge.svg)

An [Alexa skill](https://developer.amazon.com/en-US/alexa) to discover when the next tram is due using the [Traveline NextBuses API](https://www.travelinedata.org.uk/traveline-open-data/nextbuses-api/) and an [AWS Lambda](https://aws.amazon.com/lambda/) written in [Go](https://golang.org/).

## Setup

### AWS Lambda

You will require an [Amazon Web Services (AWS)](https://aws.amazon.com/account) account.

Please ensure the following dependencies are installed:

- [Git](https://git-scm.com/)
- [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- [Go](https://golang.org/) 1.14+
- [AWS CLI](https://aws.amazon.com/cli/)

Configure the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) in the WSL console.

Then run:

```shell
make install
```

To deploy the infrastructure to AWS run:

```shell
make deploy USERNAME=[] PASSWORD=[] AWS_ACCOUNT_ID=[] AWS_REGION=[] NAPTAN_CODE=[]
```

Where:

- `USERNAME` - the username for the Traveline NextBuses API
- `PASSWORD` - the password for the Traveline NextBuses API
- `AWS_ACCOUNT_ID` - the AWS account ID to install the infrastructure
- `AWS_REGION` - the AWS region to install the infrastructure
- `NAPTAN_CODE` - the NaPTAN code of the tram stop for the next tram times

### Alexa skill

You will require an [Amazon developer account](https://developer.amazon.com/).

Login to the [alexa develop console](https://developer.amazon.com/alexa/console/ask) and create a new skill with the following configuration.

#### Create a new skill

- Skill name: `Next tram`
- Default language: `English (UK)`
- Choose a model to add to your skill: `Custom`
- Choose a method to host your skill's backend resources: `Provision your own`

#### Build

Choose: `Start from scratch`

Invocation

- Skill Invocation Name: `next tram`

Intents

- Name: `when`
- Sample Utterance: `When is my next tram due`

Endpoint

- Service Endpoint Type: AWS Lambda ARN
- Default Region: _Use the ARN output when the infrastructure stack is deployed_

You will then need to build the model.

#### Test

Select the `Test` tab in teh alexa developer console.

- Skill testing is enabled in: `Development`

Alexa Simulator: `ask next tram when my next tram is due`

And the magic should happen!
