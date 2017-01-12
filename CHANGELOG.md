# Tideland Go CouchDB Client

## 2017-01-12

- Made view key parameters typing more flexible 
- View result set now also returns number of returned rows
- Added configuration option for debug logging of requests

## 2016-12-27

- Added helper to create configurations

## 2016-12-26

- Optimized retrieval of *id* and *revision*

## 2016-12-24

- Added *security* package for administrator and user management
- User authentication via basic authentication or session cookies as parameters
- Database security via `/{db}/_security` document

## 2016-12-01

- Added parameters *skip* and *limit* for views

## 2016-11-30

- Larger renamings for more consistent naming
- *View* is now an own type

## 2016-11-27

- Made parameters more beautiful
- Started better design implementation

## 2016-11-25

- Added *startup* package for database creation and
  maintenance

## 2016-11-24

- Added CRUD for design documents

## 2016-11-23

- Edit *OpenPath()* to retrieve the nested configuration for
  opening the database

## 2016-11-22

- Extended with bulk, design document, and view requests

## 2016-11-21

- Initial *alpha* with CRUD functionality
