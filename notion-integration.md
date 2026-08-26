# notion integration with vault

## Configuration

We have an integration that we can create a way for me to configure notion a specific vault. There will be a two way street where we will need to specify wehre we need to put the files in notion and where they are being synced from (what vault)

## Sync Service

We will need to create a sync service that listens on the vault (which is basiclaly just a folder) then everytime a change happens it will get what file then update it in notion

## Parser

There are packages out there that already parse form markdown into notion documents. Are we able to do that? if we cannot find any respectable ones in golang we can create this as a helper service in python   
