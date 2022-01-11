export function atob(str) {
  return Buffer.from(str, 'base64').toString('binary')
}

export function btoa(str) {
  return Buffer.from(str, 'binary').toString('base64')
}

export const removeEmpty = obj => {
  let newObj = {}
  Object.keys(obj).forEach(key => {
    if (obj[key] === Object(obj[key])) newObj[key] = removeEmpty(obj[key])
    else if (obj[key] !== undefined) newObj[key] = obj[key]
  })
  return newObj
}
