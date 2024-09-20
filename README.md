# DICEROLLER SERVER
Answer to an Interview problem. 

# How to use
This webservice has two main endpoints:
- Generate a time-based access token that is valid for three minutes
- Roll any number of any-sided dice, provided that the request has a valid access token

These endpoints are accessible over HTTP. You can use a web browser or the pre-built Postman Collection. 

Please see the below notes for each API's documentation. 

## Using Postman
1. Go to [https://www.postman.com/](https://www.postman.com/) and install Postman or use the web interface
2. Import the collection found at [Diceroller.postman_collection.json](Diceroller.postman_collection.json)
3. Generate a token from the POST endpoint by hitting the "send" button.
4. Copy the token result into the GET endpoint. Configure the dice as you see fit. You can add more!
5. "send" your die rolls!

## Using the Web
The url can be found at [direrollerserver-382425995150.us-central1.run.app](direrollerserver-382425995150.us-central1.run.app)

# Endpoints

## POST /token
This endpoint generates an access token that lives for three minutes.
There is no limit on how many tokens that you can make.

### Parameters
None

### Response
HTTP 201
```
{
    "token": "GENERATED_TOKEN"
}
```

## GET /diceroller
This endpoint wil roll any number of any sided dice, provided that the user submits a valid access token. The total die result will be returned along with a breakdown of the rolls of each die so that the GM can fudge the numbers as they see fit before announcing it to their group.

### Parameters
All parameters are Query parameters

ex: https://direrollerserver-382425995150.us-central1.run.app/diceroller?dice=1d4&dice=5d2&token=REPLACE_ME_WITH_GENERATED_TOKEN 

```
dice
    - REQUIRED
    - Query parameter
    - User dictated how many dice, and how many sides it has
    - This can be specified multiple times for multiple die sides.
    - there is no validation if the same sided die is declared multiple times.
    - case sensitive
    - ex: dice=1d4&dice=5d2
token
    - REQUIRED
    - Query parameter
    - A valid access token generated from the post endpoint
    - ex: token=ASDEsds==
```

### Response
HTTP 200
```
{
    "total": 8, // the sum of all of the die rolls
    "breakdown": { // A breakdown, by die, of each roll result
        "d2":[1,1,1,1,2], 
        "d4":[2]
    }
}
```
