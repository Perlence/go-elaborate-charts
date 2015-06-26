/// <reference path="../typings/tsd.d.ts" />
/// <reference path="../typings/app.d.ts" />
var BACKEND_URL = 'http://127.0.0.1:8080'

class App {
    alertTmpl: _.TemplateExecutor
    formTmpl: _.TemplateExecutor

    $settings: JQuery
    $form: JQuery
    $navbarToggle: JQuery
    $alerts: JQuery

    settingsForm: SettingsForm

    api: LastFMClient

    constructor() {
        this.$settings = $('#settings-block')

        this.$navbarToggle = this.$settings.find('.navbar-toggle')

        this.$alerts = this.$settings.find('#alerts')
        this.alertTmpl = _.template($('#alert-template').html())

        this.$form = this.$settings.find('#form')
        this.settingsForm = new SettingsForm(this.$form)
        this.formTmpl = _.template($('#form-template').html())

        this.api = new LastFMClient(BACKEND_URL)
    }

    bindEvents() {
        this.$form.on('submit', this.submit.bind(this))
        this.$navbarToggle.on('click', this.navbarToggle.bind(this))
        $(window).on('resize', this.resize.bind(this))
    }

    render() {
        this.settingsForm.render()
    }

    submit(e: JQueryEventObject): boolean {
        this.clearAlerts()
        var errors = this.settingsForm.verify()
        _.forEach(errors, (error: Error, input: string) => {
            this.showAlert('danger', 'Verification error', error)
        })
        if (!_.isEmpty(errors)) return false

        var uri = BACKEND_URL + '/get_weekly_chart'
        this.api.getInfo(this.settingsForm.username())
            .then((info: UserInfoResponse) => {
                var toDate = moment.utc()
                var registeredAt = moment.utc(Number(info.Registered.Unixtime) * 1000)
                var fromDate = this.settingsForm.fromDate(toDate, registeredAt)

                var ranges = spanRange(fromDate, toDate, moment.duration(1, 'week')).reverse()
                return Promise.map(ranges, (span: Span) => {
                    return this.api.getWeeklyChart(
                        this.settingsForm.username(),
                        this.settingsForm.chartType(),
                        span.start,
                        span.end
                    )
                })
            })
            .catch((err: Error) => {
                console.error(err)
                this.showAlert('danger', 'Error', err)
            })
        return false
    }

    navbarToggle(e: JQueryEventObject) {
        this.$settings.toggleClass('collapsed')
        if (this.$settings.hasClass('collapsed')) {
            this.$form.css('pointer-events', 'none')
            this.$settings.css('left', -this.$settings.width() + 72)
        }
        else {
            this.$form.css('pointer-events', 'auto')
            this.$settings.css('left', 0)
        }
    }

    resize(e: JQueryEventObject) {
        if (this.$settings.hasClass('collapsed')) {
            this.$settings.css('left', -this.$settings.width() + 72)
        }
    }

    showAlert(style: string, reason: string, err: Error) {
        var html = this.alertTmpl({
            style: style,
            reason: reason,
            message: err.message,
        })
        this.$alerts.append(html)
    }

    clearAlerts() {
        this.$alerts.empty()
    }
}

class SettingsForm {
    $form: JQuery
    $settings: JQuery
    $username: JQuery
    $chartType: JQuery
    $numberOfPositions: JQuery
    $timeframe: JQuery

    chartTypeValues = {
        album: 'Albums',
        artist: 'Artists',
        track: 'Tracks',
    }
    numberOfPositionsValues = [5, 10, 15, 20, 30, 50]
    numberOfPositionsDefault = 3
    timeframeValues = {
        'last-7-days': 'Last 7 days',
        'last-month': 'Last month',
        'last-3-months': 'Last 3 months',
        'last-6-months': 'Last 6 months',
        'last-12-months': 'Last 12 months',
        'overall': 'Overall',
    }
    timeframeDefault = 'last-6-months'

    constructor($form: JQuery) {
        this.$form = $form
        this.$username = $form.find('#username')
        this.$chartType = $form.find('#chart-type')
        this.$numberOfPositions = $form.find('#number-of-positions')
        this.$timeframe = $form.find('#timeframe')
    }

    render() {
        var option = _.template('<option value="${ value }">${ desc }</option>')
        var optionSelected = _.template('<option value="${ value }" selected>${ desc }</option>')
        _.forEach(this.chartTypeValues, (desc, value) => {
            this.$chartType.append(option({desc: desc, value: value}))
        })
        _.forEach(this.numberOfPositionsValues, (desc, value) => {
            if (value === this.numberOfPositionsDefault) {
                this.$numberOfPositions.append(optionSelected({desc: desc, value: value}))
            } else {
                this.$numberOfPositions.append(option({desc: desc, value: value}))
            }
        })
        _.forEach(this.timeframeValues, (desc, value) => {
            if (value === this.timeframeDefault) {
                this.$timeframe.append(optionSelected({desc: desc, value: value}))
            } else {
                this.$timeframe.append(option({desc: desc, value: value}))
            }
        })
    }

    username(): string { return this.$username.val() }
    chartType(): string { return this.$chartType.val() }
    numberOfPositions(): string { return this.$numberOfPositions.val() }
    timeframe(): string { return this.$timeframe.val() }

    fromDate(toDate: moment.Moment, registeredAt: moment.Moment): moment.Moment {
        var timeframe = this.timeframe()
        var fromDate: moment.Moment
        switch (timeframe) {
        case 'last-7-days':
            fromDate = toDate.clone().subtract(2, 'week')
            break
        case 'last-month':
            fromDate = toDate.clone().subtract(2, 'month')
            break
        case 'last-3-months':
            fromDate = toDate.clone().subtract(6, 'month')
            break
        case 'last-6-months':
            fromDate = toDate.clone().subtract(12, 'month')
            break
        case 'last-12-months':
            fromDate = toDate.clone().subtract(24, 'month')
            break
        case 'overall':
            fromDate = registeredAt
            break
        default:
            throw new Error('Unrecognized timeframe: ' + timeframe)
        }
        fromDate.startOf('week').add(12, 'hours')
        return fromDate
    }

    verify(): IObject<Error> {
        var errors: IObject<Error> = {}
        if (!this.username()) {
            errors['username'] = new Error('Please write your username')
        }
        if (!(this.chartType() in this.chartTypeValues)) {
            errors['chartType'] = new Error('Unrecognized chart type')
        }
        if (!(this.timeframe() in this.timeframeValues)) {
            errors['timeframe'] = new Error('Unrecognized timeframe')
        }
        return errors
    }
}

$(() => {
    var app = new App()
    app.bindEvents()
    app.render()
});
