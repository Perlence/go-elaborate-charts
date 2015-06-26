/// <reference path="../typings/tsd.d.ts" />
/// <reference path="../typings/app.d.ts" />
interface Span {
  start: moment.Moment
  end: moment.Moment
}

interface IObject<T> {
  [index: string]: T
}

function spanRange(start: moment.Moment, end: moment.Moment, duration: moment.Duration): Span[] {
  var s: moment.Moment, n: moment.Moment, e: moment.Moment
  var result: Span[] = []
  s = start.clone()
  while (s < end) {
    n = s.clone().add(duration)
    e = moment.min([end, n])
    result.push({start: s, end: e})
    s = n
  }
  return result
}
