const PeerId = require('peer-id')

const id = "Qmc8L26tfbu6hbnhruDaoj6VBS2kVAYemxBNbNVLJDJisR"

async function run() {
  const peerId = PeerId.createFromB58String(id) 
  console.log(peerId)
}
run()
