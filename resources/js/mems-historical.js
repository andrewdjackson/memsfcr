addData = function(chart, label, data, fault) {
    chart.data.labels.shift()
    chart.data.labels.push(label);
    chart.data.datasets[0].data.push(data)
    chart.data.datasets[0].data.shift()

    if (fault > 0) {
        chart.data.datasets[1].data.push(data)
        chart.data.datasets[1].borderColor = 'rgba(202,12,55,0.7)'
        chart.data.datasets[1].data.shift()
    }

    chart.update('none');
}

addScenarioData = function(chart, data) {
    chart.data = data
}

createChart = function(id, title) {
    var ctx = $('#' + id);

    return new Chart(ctx, {
        type: 'line',
        data: {
            labels: Array.apply(null, Array(120)).map(function() { return '' }),
            datasets: [{
                data: Array.apply(null, Array(120)).map(function() { return 0 }),
                cubicInterpolationMode: 'monotone',
                tension: 0.4,
                borderColor: 'rgba(102,102,255,1)',
                backgroundColor: 'rgba(102,153,204,0.1)',
                fillColor: "rgba(102,153,51,0.1)",
                strokeColor: "rgba(220,220,220,1)",
                borderWidth: 1,
                fill: true,
            },
            {
                // faults data line
                data: Array.apply(null, Array(120)).map(function() { return 0 }),
                cubicInterpolationMode: 'monotone',
                tension: 0.4,
                borderColor: 'rgba(102,102,255,0)',
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            spanGaps: true,
            radius: 0,
            plugins: {
                legend: {
                    display: false,
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    stacked: false,
                    grid: {
                        fontStyle: "normal",
                        fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                        color: "rgba(102,153,0,0.2)"
                    },
                    title: {
                        fontSize: 14,
                        fontStyle: "normal",
                        fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                        display: true,
                        text: title,
                    },
                },
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        display:false
                    },
                }
            },
        }
    });
}

createSpark = function(id) {
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
                borderWidth: 1,
                cubicInterpolationMode: 'monotone',
                tension: 0.4,
                fill: true,
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            spanGaps: true,
            radius: 0,
            plugins: {
                legend: {
                    display: false,
                }
            },
            tooltips: {
                enabled: false
            },
            scales: {
                y: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        display:false
                    },
                },
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        display:false
                    },
                }
            }
        },
    });
}

