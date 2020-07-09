addData = function (chart, label, data) {
    chart.data.labels.shift()
    chart.data.labels.push(label);
    chart.data.datasets.forEach((dataset) => {
        dataset.data.push(data);
        dataset.data.shift()
    });
    chart.update();
}

createChart = function (id, title, low, high) {
    var ctx = $('#' + id);

    return new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
            datasets: [{
                data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                borderColor: 'rgba(102,102,255,0.9)',
                backgroundColor: 'rgba(102,153,204,0.1)',
                fillColor: "rgba(102,153,51,0.2)",
                strokeColor: "rgba(220,220,220,1)",
                borderWidth: 1
            }],
        },
        options: {
            maintainAspectRatio: false,
            responsive: true,
            spanGaps: true,
            fontSize: 12,
            fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
            animation: {
                duration: 300, 
            },
            hover: {
                animationDuration: 0 
            },
            responsiveAnimationDuration: 0,
            legend: {
                display: false,
            },
            elements: {
                line: {
                    cubicInterpolationMode: 'monotone',
                }
            },
            title: {
                fontSize: 14,
                fontStyle: "normal",
                fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                display: true,
                text: title,
            },
            scales: {
                yAxes: [{
                    stacked: false,
                    gridLines: {
                        display: true,
                        fontStyle: "normal",
                        fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                        color: "rgba(102,153,0,0.2)"
                    }
                }],
                xAxes: [{
                    gridLines: {
                        display: false
                    }
                }]
            },
        },
    });
}

createSpark = function (id) {
    var ctx = $('#' + id);

    return new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
            datasets: [{
                data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                borderColor: 'rgba(102,102,255,0.9)',
                backgroundColor: 'rgba(102,153,204,0.1)',
                fillColor: "rgba(102,153,51,0.2)",
                strokeColor: "rgba(220,220,220,1)",
                borderWidth: 1
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            spanGaps: true,
            legend: {
                display: false
            },
            animation: {
                duration: 0,
            },
            hover: {
                animationDuration: 0,
            },
            responsiveAnimationDuration: 0,
            elements: {
                line: {
                    borderColor: '#000000',
                    borderWidth: 1,
                    cubicInterpolationMode: 'monotone',
                },
                point: {
                    radius: 0
                }
            },
            tooltips: {
                enabled: false
            },
            scales: {
                yAxes: [
                    {
                        display: false
                    }
                ],
                xAxes: [
                    {
                        display: false
                    }
                ]
            }
        },
    });
}