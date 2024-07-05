# rest-api-service-golang-private-1

## Database

stories
---------
id
title
content
author_id
created_at
updated_at
deleted_at

users
---------
id
username
password
created_at

## Endpoints

GET /api/v1/stories
GET /api/v1/stories/{id}
POST /api/v1/stories -> Authorization: Bearer jwttoken
PUT /api/v1/stories/{id}  -> Authorization: Bearer jwttoken
DELETE /api/v1/stories/{id} -> Authorization: Bearer jwttoken

POST /api/v1/login -> Authorization: Basic basictoken