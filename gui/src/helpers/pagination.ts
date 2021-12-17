export function paginate(arr: any[], chunk: number) {
  let result: any[][] = []
  for (let i = 0; i < arr.length; i += chunk) {
    let tempArray
    tempArray = arr.slice(i, i + chunk)
    result.push(tempArray)
  }
  return result
}

export function range(start: number, end: number): number[] {
  let length = end - start + 1
  return Array.from({ length }, (_, idx) => idx + start)
}
