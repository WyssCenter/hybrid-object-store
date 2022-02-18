# Roles
There are essentially three roles in the system. These roles are used to limit abilities for users.

## Admin Role
The `admin` role has all the abilities privileged and normal users have with additional capabilities. 

Users with this role:

* Will see and have read/write access on **all** datasets, regardless if permissions have been explicitly applied by a user. 
* Can configure namespaces and namespace syncing
* Can remove users from the system group `public`

## Privileged Role
The `privileged` role has all the abilities normal users have with additional capabilities. 

Users with this role:

* Can create and delete datasets
* Can configure dataset syncing
* Can create and manage groups

## User Role
The `user` role is the default role a user has if they do not have the `admin` or `privileged` roles. 

Users with this role:

* Interact with datasets to which they have been granted permissions
* Read/write data based on their permissions
* Create PATs

## System Groups

The system creates and maintains two groups automatically. 

The `admin` system group will contain all users with the `admin` role. This group is automatically attached to all datasets at creation with read-write 
access. You cannot remove this group from a dataset.

The `public` system group will contain all users in the system. This group is automatically populated when users log in. It can be used by anyone in the
group to grant read-only access to datasets. You cannot apply read-write permissions with this group.
