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
| /login | POST | Authenticates user | |
| /user/entries/:id | GET | JSON of all items | |
| /user/entry/create/:id | POST | Create new item return JSON of new item | |
| /user/entry/update/:id/:medication_id | PUT | update item with matching id, return its JSON |  |
| /user/entry/delete/:id/:medication_id | DELETE | delete the item with the matching id | |
| /user/logout | POST | Logout user | |



## ERD
![Entity Relationship Diagram](https://i.imgur.com/0Gxp1Cy.png)
