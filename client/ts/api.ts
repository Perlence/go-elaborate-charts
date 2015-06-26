/// <reference path="../typings/tsd.d.ts" />
/// <reference path="../typings/app.d.ts" />
interface UserInfoResponse {
    XMLName: string
    Id: string
    Name: string
    RealName: string
    Url: string
    Country: string
    Age: string
    Gender: string
    Subscriber: string
    PlayCount: string
    Playlists: string
    Bootstrap: string
    Registered: {
        Unixtime: string
        Time: string
    }
    Type: string
    Images: {
        Size: string
        Url: string
    }[]
}

interface WeeklyChartResponse {
    Chart: Object
    ToDate: number
}

interface ErrorResponse {
    Error: string
}

class LastFMClient {
    private baseUrl: string

    constructor(baseUrl) {
        this.baseUrl = baseUrl
    }

    private getJSON<T>(url: string, data?: Object|string): Promise<T> {
        return new Promise<T>((resolve, reject) => {
            $.getJSON(this.baseUrl + url, data)
                .done((result) => {
                    resolve(result)
                })
                .fail((jqxhr, textStatus, reason) => {
                    var response: ErrorResponse = jqxhr.responseJSON
                    var message = (response != null) ? response.Error : reason
                    reject(new XHRError(message, jqxhr, textStatus, reason))
                })
        })
    }

    getInfo(username: string): Promise<UserInfoResponse> {
        return this.getJSON<UserInfoResponse>('/get_info', {username: username})
    }

    getWeeklyChart(username: string, chartType: string, fromDate: moment.Moment, toDate: moment.Moment): Promise<WeeklyChartResponse> {
        return this.getJSON<WeeklyChartResponse>('/get_weekly_chart', {
            username: username,
            chart_type: chartType,
            from_date: fromDate.unix(),
            to_date: toDate.unix()
        })
    }
}
