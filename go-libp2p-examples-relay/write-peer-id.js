const fs = require('fs')
const PeerId = require('peer-id')
const obj = require('./peer-id.json')

async function run () {
  const peerId = await PeerId.createFromJSON(obj)
  console.log('Jim', peerId)
  const peerIdBin = Buffer.from(peerId.marshalPrivKey())
  console.log('Jim2', peerIdBin)
  fs.writeFileSync('privkey', peerIdBin.toString('base64'))
}
run()
