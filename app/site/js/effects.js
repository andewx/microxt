//Effectual functions with jQuery and CSS3

//Pulse Element Text Color
function pulseTextColor(element, color1, color2, duration) {
    var factor = 0.0;
    var interval = setInterval(function () {
        var factor = 0.0;
        var color = r2h(_interpolateColor(h2r(color1), h2r(color2), factor));
        element.css("color", color);
        factor += 0.1;
        if (factor > 1.0) {
            factor = 0.0;
        }
    }, duration);


    //Interpolate between color one and color two
    return interval;
}


function pulseBackgroundColor(element, color1, color2, duration) {
    var color = color1;
    var interval = setInterval(function () {
        if (color == color1) {
            color = color2;
        } else {
            color = color1;
        }
        element.css("background-color", color);
    }, duration);
}


function textCarousel(element, str, duration) {
    var index = 0;
    var interval = setInterval(function () {
        var rindex = Math.floor(Math.random() * str.length);
        element.text(str[rindex]);
        index++;
        if (index == str.length) {
            index = 0;
        }
    }, duration);
    var size = 0.5;
    var step_size = 0.1;
    var size_interval = setInterval(function () {
        size += step_size;
        element.css("font-size", size + "em");
        if (size > 1.5) {
            step_size = -0.1;
        }
        if (size < 0.5) {
            step_size = 0.1;
        }
    }, duration/10);
}


/**
 * TODO - refactor this as a (jQuery?) plugin!
**/

// Converts a #ffffff hex string into an [r,g,b] array
var h2r = function(hex) {
    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? [
        parseInt(result[1], 16),
        parseInt(result[2], 16),
        parseInt(result[3], 16)
    ] : null;
};

// Inverse of the above
var r2h = function(rgb) {
    return "#" + ((1 << 24) + (rgb[0] << 16) + (rgb[1] << 8) + rgb[2]).toString(16).slice(1);
};

// Interpolates two [r,g,b] colors and returns an [r,g,b] of the result
// Taken from the awesome ROT.js roguelike dev library at
// https://github.com/ondras/rot.js
var _interpolateColor = function(color1, color2, factor) {
  if (arguments.length < 3) { factor = 0.5; }
  var result = color1.slice();
  for (var i=0;i<3;i++) {
    result[i] = Math.round(result[i] + factor*(color2[i]-color1[i]));
  }
  return result;
};

var rgb2hsl = function(color) {
  var r = color[0]/255;
  var g = color[1]/255;
  var b = color[2]/255;

  var max = Math.max(r, g, b), min = Math.min(r, g, b);
  var h, s, l = (max + min) / 2;

  if (max == min) {
    h = s = 0; // achromatic
  } else {
    var d = max - min;
    s = (l > 0.5 ? d / (2 - max - min) : d / (max + min));
    switch(max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return [h, s, l];
};

var hsl2rgb = function(color) {
  var l = color[2];

  if (color[1] == 0) {
    l = Math.round(l*255);
    return [l, l, l];
  } else {
    function hue2rgb(p, q, t) {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1/6) return p + (q - p) * 6 * t;
      if (t < 1/2) return q;
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
      return p;
    }

    var s = color[1];
    var q = (l < 0.5 ? l * (1 + s) : l + s - l * s);
    var p = 2 * l - q;
    var r = hue2rgb(p, q, color[0] + 1/3);
    var g = hue2rgb(p, q, color[0]);
    var b = hue2rgb(p, q, color[0] - 1/3);
    return [Math.round(r*255), Math.round(g*255), Math.round(b*255)];
  }
};

var _interpolateHSL = function(color1, color2, factor) {
  if (arguments.length < 3) { factor = 0.5; }
  var hsl1 = rgb2hsl(color1);
  var hsl2 = rgb2hsl(color2);
  for (var i=0;i<3;i++) {
    hsl1[i] += factor*(hsl2[i]-hsl1[i]);
  }
  return hsl2rgb(hsl1);
};


(function($) {
  
  var $list = $('#list'),
      $start = $('#start'),
      $end = $('#end'),
      $intype = $('input[name="intype"]'),
      $usehex = $('#usehex');

  // Add li elements between the start and end ones
  var _createSteps = function(numSteps) {
    $list.find('li.interim').remove();

    for(var i = 0; i < numSteps; i++) {
      $end.before('<li class="interim"><span></span></li>');
    }
  };

  // Color each li by interpolating between the start and end colors
  var _styleSteps = function() {
    var $items = $('li'),
        scol = h2r($start.find('input').val()),
        ecol = h2r($end.find('input').val());
    var fn = '_' + $('input[name="intype"]:checked').val();
    
    console.log('fn', fn);

    var factorStep = 1 / ($items.length - 1);
    $items.each(function(idx) {
      var icol = window[fn](scol, ecol, factorStep * idx),
          hcol = r2h(icol);

      $(this).css('background-color', hcol);
      $(this).find('span').text(hcol);
    });

  }

  // Re-render on change
  $('#usehex').on('change', function() {
    var ct = $('ul input').eq(0).attr('type');
    var scol = $start.find('input').val(),
        ecol = $end.find('input').val();
    
    $('ul input').attr('type', (ct == 'color') ? 'text' : 'color');
    
    $start.find('input').val(scol);
    $end.find('input').val(ecol);
  });
  $('input').not('#usehex').on('change', _styleSteps);
  $('#numsteps').on('change', function() {
    _createSteps($(this).val());
    _styleSteps();
  }).trigger('change');
  
})(jQuery);



