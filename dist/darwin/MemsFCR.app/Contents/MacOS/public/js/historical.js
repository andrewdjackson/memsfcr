$(document).ready(function() {
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
                    label: title,
                    data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                    borderColor: 'rgba(192, 192, 192, 1)',
                    backgroundColor: 'rgba(192, 192, 192, 0.2)',
                    fillColor: "rgba(0,153,204,0.2)",
                    strokeColor: "rgba(220,220,220,1)",
                    borderWidth: 1
                }],
                yHighlightRange: {
                    begin: low,
                    end: high
                }
            },
            options: {
                maintainAspectRatio: false,
                scales: {
                    yAxes: [{
                        stacked: true,
                        gridLines: {
                            display: true,
                            color: "rgba(75, 192, 192, 0.2)"
                        }
                    }],
                    xAxes: [{
                        gridLines: {
                            display: false
                        }
                    }]
                },
                options: {
                    fontFamily: "'Segoe UI', 'Tahoma', 'Geneva', 'Verdana', sans-serif",
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
});