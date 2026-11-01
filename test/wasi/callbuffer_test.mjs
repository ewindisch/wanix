import assert from "node:assert/strict";
import { decode } from "cbor-x";
import { CallBuffer } from "../../wasi/worker/lib.js";

const buffer = new SharedArrayBuffer(16384);
const call = new CallBuffer(buffer);
const ctrl = new Int32Array(buffer, 0, 2);
const len = new Int32Array(buffer, 4, 1);
const data = new Uint8Array(buffer, 8);
const response = () => decode(data.slice(0, len[0]));

call.respond(new Uint8Array(call.maxData));
assert.equal(ctrl[0], 1);
assert.equal(response().byteLength, call.maxData);

Atomics.store(ctrl, 0, 0);
call.respond(new Uint8Array(call.maxData + 1));
assert.equal(ctrl[0], 1);
assert.deepEqual(response(), { __wanix_error: `CallBuffer response too large: ${call.maxData + 1} bytes exceeds ${call.maxData}-byte limit` });

Atomics.store(ctrl, 0, 0);
call.respond({ value: 'ok' });
assert.deepEqual(response(), { value: 'ok' });
console.log('CallBuffer response-size tests passed');
