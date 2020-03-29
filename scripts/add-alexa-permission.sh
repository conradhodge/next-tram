#!/bin/bash

LAMBDA_ARN=$(aws lambda get-function-configuration --function-name=get-next-tram-lambda | jq .FunctionArn | tr -d '"')

printf "\nRemoving permission to allow Alexa to access the lambda...\n"

aws lambda remove-permission \
    --function-name ${LAMBDA_ARN} \
    --statement-id alexa \
    --output text

printf "\nAdding permission to allow Alexa to access the lambda...\n"

aws lambda add-permission \
    --function-name ${LAMBDA_ARN} \
    --action lambda:InvokeFunction \
    --statement-id alexa \
    --principal alexa-appkit.amazon.com \
    --output text
