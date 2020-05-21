addData = function(chart, label, data) {
    chart.data.labels.shift()
    chart.data.labels.push(label);
    chart.data.datasets.forEach((dataset) => {
        dataset.data.push(data);
        dataset.data.shift()
    });
    chart.update();
}

createChart = function(id, title, low, high) {
    var ctx = $('#' + id);

    // The original draw function for the line chart. This will be applied after we have drawn our highlight range (as a rectangle behind the line chart).
    var originalLineDraw = Chart.controllers.line.prototype.draw;
    // Extend the line chart, in order to override the draw function.
    Chart.helpers.extend(Chart.controllers.line.prototype, {
        draw: function() {
            var chart = this.chart;
            // Get the object that determines the region to highlight.
            var yHighlightRange = chart.config.data.yHighlightRange;

            // If the object exists.
            if (yHighlightRange !== undefined) {
                if (yHighlightRange.begin !== undefined) {
                    var ctx = chart.chart.ctx;

                    var yRangeBegin = yHighlightRange.begin;
                    var yRangeEnd = yHighlightRange.end;

                    var xaxis = chart.scales['x-axis-0'];
                    var yaxis = chart.scales['y-axis-0'];

                    var yRangeBeginPixel = yaxis.getPixelForValue(yRangeBegin);
                    var yRangeEndPixel = yaxis.getPixelForValue(yRangeEnd);

                    if (yaxis.max > yRangeBegin) {
                        ctx.save();

                        // The fill style of the rectangle we are about to fill.
                        ctx.fillStyle = 'rgba(127, 191, 63, 0.05)';
                        // Fill the rectangle that represents the highlight region. The parameters are the closest-to-starting-point pixel's x-coordinate,
                        // the closest-to-starting-point pixel's y-coordinate, the width of the rectangle in pixels, and the height of the rectangle in pixels, respectively.
                        ctx.fillRect(xaxis.left, Math.min(yRangeBeginPixel, yRangeEndPixel), xaxis.right - xaxis.left, Math.max(yRangeBeginPixel, yRangeEndPixel) - Math.min(yRangeBeginPixel, yRangeEndPixel));

                        ctx.restore();
                    }
                }
            }

            // Apply the original draw function for the line chart.
            originalLineDraw.apply(this, arguments);
        }
    });

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
            yHighlightRange: {
                begin: low,
                end: high
            }
        },
        options: {
            legend: {
                display: false,
            },
            title: {
                fontSize: 14,
                fontStyle: "normal",
                fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                display: true,
                text: title,
            },
            maintainAspectRatio: false,
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
            options: {
                fontSize: 12,
                fontFamily: "'-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
                animation: {
                    duration: 100 // general animation time
                },
                hover: {
                    animationDuration: 0 // duration of animations when hovering an item
                },
                responsiveAnimationDuration: 0 // animation duration after a resize
            },
        },
    });
}