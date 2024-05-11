# Capstone

# Product Requirements Documentation
| Field | Detail |
|-------|--------|
| Project Name | MediManage |
| Description | A medication management application for improved health, reduced side effects, fewer medication errors, and increased sense of control over health. |
| Developers | Timothy Rodriguez |
| Live Website | https://main--medimanage.netlify.app/ |
| Repo | https://github.com/timorodr/go-react-final-BE |
| Technologies | Golang, MongoDB |

## Problem Being Solved and Target Market

Medication management is a crucial aspect of maintaining good health, independence, and well-being, for people who rely on medication to manage chronic conditions or other health concerns.

## User Stories


- Users should be able to see the site on desktop and mobile
- Users can create an account
- Users can sign in to their account
- Users can create a new item
- Users can see all their items on the dashboard
- Users can update items
- User can delete items

## Route Tables

List of all routes and their functionality in the app

| Endpoint | Method | Response | Other |
| -------- | ------ | -------- | ----- |
| /signup | POST | Creates a new user | |
| /login | GET | Retrieves user information | |
| /login | POST | Authenticates user | |
| /medications | GET | JSON of all items | |
| /medications | POST | Create new item return JSON of new item | |
| /medications/:id | GET | JSON of item with matching id number | |
| /medications/:id | PUT | update item with matching id, return its JSON |  |
| /medications/:id | DELETE | delete the item with the matching id | |
| /interactions | POST | Checks for drug interactions with medication names in request body | |



## ERD
![Entity Relationship Diagram](https://i.imgur.com/0Gxp1Cy.png)
