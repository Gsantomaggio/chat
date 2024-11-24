function isSameSocket(s1, s2) {
  return (
    s1.remoteAddress === s2.remoteAddress && s1.remotePort === s2.remotePort
  );
}
function getSocketsExcluding(sockets, sock) {
  return sockets.filter((s) => !isSameSocket(s, sock));
}
function socketToId(sock) {
  return `${sock.remoteAddress}:${sock.remotePort}`;
}
export { getSocketsExcluding, socketToId };
