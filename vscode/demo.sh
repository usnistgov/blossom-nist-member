#!/bin/sh

export MICROFAB_CONFIG='{
    "port": 8080,
    "endorsing_organizations":[
        {
            "name": "Blossom"
        },
        {
            "name": "A1"
        },
        {
            "name": "A2"
        }
    ],
    "channels":[
        {
            "name": "channel1",
            "endorsing_organizations":[
                "Blossom",
                "A1",
                "A2"
            ],
            "capability_level": "V1_4_2"
        }
    ]
}'

docker run -p 8080:8080 -e MICROFAB_CONFIG ibmcom/ibp-microfab -e FABRIC_LOGGING_SPEC='DEBUG'
