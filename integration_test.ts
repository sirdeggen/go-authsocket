import { AuthSocketClient } from '@bsv/authsocket-client'
import { PrivateKey } from '@bsv/sdk'

// Start the Go authsocket server on ws://localhost:8080 before running this test

async function runIntegrationTest() {
  // Alice's private key (hex)
  const alicePrivHex = '0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20'
  const alicePriv = PrivateKey.fromHex(alicePrivHex)

  // Bob's private key (different)
  const bobPrivHex = '02030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f21'
  const bobPriv = PrivateKey.fromHex(bobPrivHex)

  // Alice connects
  const alice = new AuthSocketClient('ws://localhost:8080', alicePriv)
  alice.on('message', (data: any) => {
    console.log('Alice received:', data)
  })

  // Bob connects
  const bob = new AuthSocketClient('ws://localhost:8080', bobPriv)
  bob.on('message', (data: any) => {
    console.log('Bob received:', data)
  })

  // Wait for connections
  await new Promise(resolve => setTimeout(resolve, 1000))

  // Alice sends a message (this will be broadcasted by the Go server)
  alice.emit('message', { from: 'Alice', text: 'Hello from Alice!' })

  // Bob sends a message
  bob.emit('message', { from: 'Bob', text: 'Hello from Bob!' })

  // Wait for messages
  await new Promise(resolve => setTimeout(resolve, 2000))

  console.log('Integration test completed')
}

// Run the test
runIntegrationTest().catch(console.error)