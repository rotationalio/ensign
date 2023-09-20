/*
Package south provides a mechanism for database migrations to be run at startup rather
than manually running migrations using the command line. The current migration of the
data is stored in the database and if there is a migration that is later than the
current version of the data, it is applied to the database one migration at a time.
*/
package south
