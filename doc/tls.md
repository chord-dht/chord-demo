# TLS

## Scenario and TLS Principles

Currently, there are multiple peers that can act as both servers and clients to communicate with other peers. The following outlines the process of peers connecting and continuing communication using TLS:

1. **TLS Connection Process**:
   - **Generate Private Key and CSR**: Each peer generates its own private key and certificate signing request (CSR).
   - **Request Signed Certificate**: The CSR is sent to the CA server, which verifies it and returns a signed certificate.
   - **Establish Connection**:
     1. **Handshake Start**: The client peer sends a ClientHello message to the server peer, which includes supported encryption algorithms and a random number.
     2. **Server Response**: The server peer returns a ServerHello message, selecting the encryption algorithm and sending the server certificate and a random number.
     3. **Certificate Verification**: The client peer verifies whether the server's certificate is issued by a trusted CA.
     4. **Generate Session Key**: The client peer generates a pre-master secret, encrypts it with the server's public key, and sends it to the server.
     5. **Server Decryption**: The server peer decrypts the pre-master secret using its own private key to generate the session key.
     6. **Handshake Completion**: Both parties use the session key to encrypt subsequent communication, completing the handshake.

2. **Subsequent TLS Connection Process**:
   - **Encrypted Communication**: After the handshake, both parties use a symmetric encryption algorithm and the session key for encrypted communication to ensure data confidentiality and integrity.
   - **Session Management**: Session resumption mechanisms (such as session tickets or session IDs) can be used to optimize subsequent connections and reduce handshake overhead.

### How a Client Verifies the Server Certificate is Issued by a Trusted CA

Each client has built-in root certificates from CAs to verify whether the public key certificate from the server is issued by a trusted CA. The process includes:

1. **Server Sends Certificate**: When the client establishes a connection with the server, the server sends its public key certificate to the client.
2. **Client Checks Certificate Chain**: The client checks the signature chain of the server certificate to ensure each level is signed by the preceding certificate, up to the root certificate.
3. **Root Certificate Verification**: The client uses the built-in root certificate to verify the server certificate's signature. If the server certificate is signed by a trusted root or intermediate certificate, the verification passes.
4. **Certificate Validity Check**: The client also checks the certificate's validity period and revocation status to ensure it is valid and not revoked.

### Does the Client Need to Ask the CA Server Every Time?

- **No Need to Ask Every Time**: The client does not need to query the CA server every time. It uses locally stored root certificates to verify the server certificate's signature.
- **Certificate Revocation Check**: The client may use the Online Certificate Status Protocol (OCSP) or Certificate Revocation Lists (CRLs) to check if the certificate is revoked. This may involve communication with the CA server or other designated server, but it does not need to be done for every connection. It depends on the client's configuration and policy.

In summary, the client mainly relies on locally stored root certificates to verify the server certificate's signature and does not need to query the CA server every time, but certificate revocation checks may still be performed.

## Why TLS Communication is Needed for Every Process

What if TLS encrypted communication is only considered for the file transfer process?

Definition of Man-in-the-Middle (MitM):

- Has the Chord software with built-in CA root certificates.
- Obtained a public key signed certificate from the CA.

Scenario 1:

- Peer A and peer B transfer files using TLS communication.
- A man-in-the-middle also has a trusted CA root certificate and requests a public key signed certificate from the CA.
- The man-in-the-middle intercepts the public key signed certificate from server peer A and replaces it with its own.
- However, since the IP included in the certificate does not match the IP the client peer B was trying to connect to, the client peer discovers the issue, leading to a TLS connection failure, rendering the attack ineffective.

Scenario 2:

Unfortunately, in the chord ring, the IP information about other nodes held by the client is obtained from other nodes. Only its own IP and the IP of the initially selected node to join can be trusted; other IPs are acquired through communication, such as after a find successor operation.

If a man-in-the-middle intercepts the communication from your chosen joining node and replaces the found node's IP information with its own;

The lookup and obtained results are not protected by TLS, allowing the client to connect to the man-in-the-middle;

Simultaneously, the man-in-the-middle intercepts and sends the notify message intended for the found successor from the client to the original successor;

The result is that the man-in-the-middle replaces the client in joining the ring.

---

Therefore, TLS communication is needed for all processes. As long as the node you initially choose to join is not a man-in-the-middle (trusted), you don't have to worry about a MitM attack.

## How to Implement

Each peer, as a client, must have a built-in CA root certificate to verify if a server's certificate is issued by a trusted CA.

Each peer, as a server, must request a public key signed certificate from the CA. Generally, this process requires domain verification or IP verification to confirm your identity.

Therefore, each peer initially includes a CA root certificate, and ultimately needs a (private key and) public key signed certificate.

## However, in our Experiment

Our AWS servers do not have a domain, so Let's Encrypt cannot be used.

Although there is an IP, applying for and obtaining a CA-signed public key signed certificate for each server is too complex, so it was not adopted.

Therefore, we chose to set up our own CA for a temporary solution in the experiment, involving pre-shared keys, specifically:

1. Generate a CA private key and self-signed root certificate.
2. Create an OpenSSL configuration file (e.g., `openssl.cnf`) to define CA settings.
3. **Pre-generate a private key and CSR for each peer, using the CA to sign the peer's certificate**.
4. **Finally, at the time of use, distribute the certificates and keys**. Each peer holds three items: the CA root certificate, a private key, and a public key signed certificate.

The reason for this setup is because our self-established CA is not a server, so it cannot handle public key signed certificate requests from each peer afterward or perform certificate revocation checks for peers as clients. Thus, this distribution is used as a last resort.

Additionally, since our certificates do not contain IP information due to the issuance process without entries, we need to enable `InsecureSkipVerify: true` to bypass the normal verification process, only using the CA root certificate to verify the server's public key signed certificate.
