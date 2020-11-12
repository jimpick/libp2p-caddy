const fs = require('fs')
const PeerId = require('peer-id')

const privKeyBase64 = fs.readFileSync('privkey', 'utf8')

async function run () {
  console.log('Jim privKeyBase64', privKeyBase64)
  const peerId = await PeerId.createFromPrivKey(Buffer.from(privKeyBase64, 'base64'))
  console.log('Jim', peerId)
}
run()
