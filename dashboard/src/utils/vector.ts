import { IVector3 } from '../types/3d'
export const adjustY = (vector: IVector3, y: number): IVector3 => {
  return [
    vector[0],
    y,
    vector[2],
  ]
}