# HomeGate -- Your Automatic Gate to your Home 

This repository contains WIP code for controller an electronic gate controller by RF. 
It uses a Raspberry Pi as a transmitter to the gate. Opening a gate can be triggered 
via a command to an *HomeGate* server.

# Build the Frontend

Make sure npm is installed:

```
    cd server/cmd/homegate/frontend
    npm install
    # to serve for testing
    # npm run serve

    # to build dist for production
    npm build
```

## Credit
- [Yaron Weinsberg](wyaron@gmail.com)
- [Tal Anker](tal.anker@gmail.com)
