export function atob(str) {
  return Buffer.from(str, 'base64').toString('binary')
}

export function btoa(str) {
  return Buffer.from(str, 'binary').toString('base64')
}
