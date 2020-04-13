#!/bin/bash

FUNCTION_NAME=get-next-tram-lambda
STATEMENT_ID=alexa

if ! OUTPUT="$(aws lambda get-policy --function-name ${FUNCTION_NAME} 2>&1)"; then
    if [ $(echo ${OUTPUT} | grep -c "ResourceNotFoundException") -gt 0 ]
    then
        printf "Adding resource-based policy to allow Alexa to access the lambda...\n"

        aws lambda add-permission \
            --function-name ${FUNCTION_NAME} \
            --statement-id ${STATEMENT_ID} \
            --principal alexa-appkit.amazon.com \
            --action lambda:InvokeFunction \
            --output text
    else
        echo "Unexpected error:\n${OUTPUT}"
    fi
fi
