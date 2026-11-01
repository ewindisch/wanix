import assert from "node:assert/strict";
import { File, OpenFile } from "../../wasi/fs.ts";

let sizeQueries = 0;
const handle = {
  open() {}, close() {}, flush() {},
  read() { return 0; }, write() { return 0; }, truncate() {},
  getSize() { sizeQueries++; throw new Error("path is unlinked"); },
};

const fd = new OpenFile(new File(handle));
const { ret, fdstat } = fd.fd_fdstat_get();
assert.equal(ret, 0);
assert.equal(fdstat.fs_filetype, 4); // WASI regular file
assert.equal(sizeQueries, 0);
console.log("fd_fdstat_get does not stat an open file by path");
