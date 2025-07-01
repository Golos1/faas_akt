# Purpose

Package of Goakt actors which take advantage of FaaS offerings on cloud providers, mostly by wrapping the clients in their SDK in some way. Actors that directly invoke specific functions will return their results as JSON, while those thand send events will just confirm success sending the event.

## Currently implemented

    1. AWS Lambda Actor (Returns JSON results). Package: lambda
    2. Inngest Event Sender (Does not return results, only confirms event was sent.) Package: inngest.
