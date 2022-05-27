'use strict'

var express = require('express')
var bodyParser = require('body-parser')

var app = express()
app.use(bodyParser.json())

const { Gateway, Wallets } = require('fabric-network')
const path = require('path')
const fs = require('fs')
const { logger } = require('./log')
const { exit } = require('process')

async function GatewayConnect(){

    const ccpPath_org1 = path.resolve(__dirname, '..', 'network', 'organizations', 'peerOrganizations', 'org1.blin.com', 'connection-org1.json')

    const ccpPath_org2 = path.resolve(__dirname, '..', 'network', 'organizations', 'peerOrganizations', 'org2.blin.com', 'connection-org2.json')

    const ccp_org1 = JSON.parse(fs.readFileSync(ccpPath_org1, 'utf8'))

    const ccp_org2 = JSON.parse(fs.readFileSync(ccpPath_org2, 'utf8'))

    const walletPath = path.join(process.cwd(), 'wallet')
    const wallet = await Wallets.newFileSystemWallet(walletPath)

    const identity = await wallet.get('appUser')

    if (!identity) {

        logger.error('appUser 유저의 ID가 지갑에 존재하지 않습니다.\nregisterUser.js를 실행하면 해결할 수 있습니다.')
        exit()

    }

    const gateway_org1 = new Gateway()

    await gateway_org1.connect(ccp_org1, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } })

const gateway_org2 = new Gateway()

    await gateway_org2.connect(ccp_org2, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } })

    global.gateway_org1 = gateway_org1
    global.gateway_org2 = gateway_org2

    return

}

app.post('/api/create_token', async function (req, res) {
    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.submitTransaction('CreateToken', req.body.name, req.body.symbol, req.body.totalsupply)

        logger.info(`POST /api/create_token : ${result.toString()}`)

        res.status(200).json({response: "200 ok"})
        
        return

    } catch (error) {

        logger.error(`POST /api/create_token error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }
})

app.get('/api/get_token_info', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Token_Info')

        logger.info(`POST /api/get_token_info : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_token_info error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/create_account', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.submitTransaction('CreateAccount', req.body.id, req.body.amount)
        logger.info(`POST /api/create_account : ${result.toString()}`)

        res.status(200).json({"response":result.toString()})
        
        return

    } catch (error) {

        logger.error(`POST /api/create_account error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/get_account', async function (req, res) {

    try {

        const network = await gateway_org2.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Account', req.body.id)

        logger.info(`POST /api/get_account : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_account error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/transfer', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.submitTransaction('Transfer', req.body.from, req.body.to, req.body._value)

        logger.info(`POST /api/transfer : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/transfer error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})


app.post('/api/get_tx', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_tx', req.body.value) //
        logger.info(`POST /api/get_tx : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_tx error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/get_root_receipt', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Root_Receipt', req.body.value)

        logger.info(`POST /api/get_root_receipt : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_root_receipt error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/get_last_receipt', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Last_Receipt', req.body.value)
    
        logger.info(`POST /api/get_last_receipt : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_last_receipt error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/get_receipts', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Receipt', req.body.value)

        logger.info(`POST /api/get_receipts : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_receipts error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.post('/api/get_receipt', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Receipts', req.body.value)

        logger.info(`POST /api/get_receipt : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`POST /api/get_receipt error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

app.get('/api/get_transfer', async function (req, res) {

    try {

        const network = await gateway_org1.getNetwork('mychannel')
        const contract = network.getContract('basic')

        const result = await contract.evaluateTransaction('Get_Transfer')

        logger.info(`GET /api/get_transfer : ${result.toString()}`)

        res.status(200).json(JSON.parse(result.toString()))
        
        return

    } catch (error) {

        logger.error(`GET /api/get_transfer error msg :  ${error}`)

        res.status(404).json({"response":"400"})
    }

})

// app.post('/api/client_register', async function (req, res) {
//     try {

//         const client_id = req.body.id
//         const client_secret = req.body.password

//         if(!client_id || !client_secret){
//             res.status(404).json({"response":"client_id or client_secret is null"})
//         }

//         // load the network configuration
//         const ccpPath = path.resolve(__dirname, '..', 'network', 'organizations', 'peerOrganizations', 'org1.blin.com', 'connection-org1.json')
//         const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'))

//         // Create a new CA client for interacting with the CA.
//         const caURL = ccp.certificateAuthorities['ca.org1.blin.com'].url
//         const ca = new FabricCAServices(caURL)

//         // Create a new file system based wallet for managing identities.
//         const walletPath = path.join(process.cwd(), 'wallet')
//         const wallet = await Wallets.newFileSystemWallet(walletPath)


//         // Check to see if we've already enrolled the user.
//         const userIdentity = await wallet.get('appUser')
//         if (userIdentity) {
//             // console.log('An identity for the user "appUser" already exists in the wallet')
//             return
//         }

//         // Check to see if we've already enrolled the admin user.
//         const adminIdentity = await wallet.get('admin')
//         if (!adminIdentity) {
//             // console.log('An identity for the admin user "admin" does not exist in the wallet')
//             // console.log('Run the enrollAdmin.js application before retrying')
//             return
//         }

//         // build a user object for authenticating with the CA
//         const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type)
//         const adminUser = await provider.getUserContext(adminIdentity, 'admin')

//         // Register the user, enroll the user, and import the new identity into the wallet.
//         const secret = await ca.register({
//             affiliation: 'org1.department1',
//             enrollmentID: 'appUser',
//             role: 'client'
//         }, adminUser)

//         const enrollment = await ca.enroll({
//             enrollmentID: 'appUser',
//             enrollmentSecret: secret
//         })
//         const x509Identity = {
//             credentials: {
//                 certificate: enrollment.certificate,
//                 privateKey: enrollment.key.toBytes(),
//             },
//             mspId: 'Org1MSP',
//             type: 'X.509',
//         }
//         await wallet.put('appUser', x509Identity)
//         // console.log('Successfully registered and enrolled admin user "appUser" and imported it into the wallet')

//     } catch (error) {
//         // console.error(`Failed to register user "appUser": ${error}`)
//         // process.exit(1)
//     }
// })

logger.info('Running on API Server http://localhost:3000')

app.listen(3000, '0.0.0.0', async () => { await GatewayConnect()})
// console.log('Running on api server')
