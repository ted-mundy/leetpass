# leetpass API

This is an incredibly simple API with a few primary functions.

## Features

### Data Signing

In order to ensure that bad actors do not misuse the system and use the app to
"airdrop" (in essence) _inappropriate photos_, all clients should only trust
data payloads that have been signed by the API. Ideally, this step wouldn't be
needed, however it is just the nature of the beast that if a system exists,
people will misuse it if the opportunity is there to do so.

The data signing/verification step will be the "authorization" part of the app,
where clients will request every *𝑥* minutes, where *𝑥* is some server-side
variable included in the response (JWT's `exp` key). Clients tokens will expire
so that if for whatever reason we decide that some of their data is 
invalid/bannable, we will no longer sign their data, thus "banning" them from 
other users.

Clients will make requests to this endpoint when _either_:

- Their tokens expire. In practice, it will be best to request this a bit sooner
than the expiration time (~5 minutes), in case of network outages or similar.
- Their broadcast data changes - username change, photo change or anything that
other people see. The server will then validate these _before_ returning the
signed data, so that for the vast majority of updates, clients receive signed
data that is OK to be broadcast, allowing for offline data access as the validation
has already been performed.
For example: Alice updates her tagline to _"I love cats!"_. The validation
function ensures that the updated data does not contain profanities, or breaks
any of the other rules. If it passes successfully, Alice receives her signed
data, which John can receive whilst offline - confident that there is nothing
which breaks the rules.

For retroactive action, something slipping through the net, the offending data
will be broadcast for *𝑥* minutes, after which point the data will be invalid,
and the server will no longer sign the data. This is obviously not ideal, but
there is no other way to do this without more API calls which - whilst possible -
isn't in scope at the moment. Maybe in the future!

There will be two endpoints for this feature as such, one to validate the data,
and one to receive the server's public key. The JWTs will be signed with RSA
encryption, so that data validation can be performed client-side.

As one of the key priorities of the project is data privacy, the API stores the
following values for each data signing request:

- Account UUID: Randomly generated unique (_the odds of collision are 1 in 2.71 x 10^18_)
account identifier. Used to manage suspended accounts.

### Image Uploading

The app has a "Photo Wall", where passersby can set an account-wide photo to be
shared and set as their designated photo to represent themselves. As a result
of the rules laid in [Data Signing](#data-signing), we cannot allow these
images to be shared P2P, and must be uploaded to our server. When an image is
uploaded, it is sent for validation _before returning a response_. The client
will receive an image ID, rather than a full URL. This is what should be
broadcast, and the recieving client should construct the image URL with the ID
and the configured URL prefix.

To avoid people using this as free image storing, the API will receive the
_Device UUID_ and _Account UUID_ of the requesting device, and will only allow
one image uploaded per device UUID and account UUID. Both aren't validated by
us, though, which means that if someone really wanted to use our object store
as file storage, they probably could. We will have some incredibly basic checks
in place here though, such as logging IPs, etc. But we are powerless really.

If you upload an image, we have to store it on our servers for validation. We
cannot allow people to go around sending images P2P, that would bypass the
rules we have laid out here. This image will be tied to the account UUID,
which is anonymous to us, as we don't have a way to link any of your account
data (username, tagline, etc) to your account UUID - we don't keep the logs!

### Encounter Ping

> [!NOTE]  
> This is by-default opt-in, but _100% anonymized_.

When an encounter happens, the API will get notified, _purely for statistics_.
When an encounter happens, _both devices_ will send some payload that contains
some randomly generated _encounter ID_.

When the server receives the encounter notification, it will do the following:

```
encounterExists := doesEncounterExist(data.encounterId)
if (encounterExists) markEncounterHandshakeComplete(data.encounterId)
else startEncounterHandshake(data.encounterId)
```

The _handshake_ is to ensure that both sides of the party have opted in to the
encounter pinging. If no response is received within a sensible time frame, it
can be assumed the other user has opted out.
