//Radar component effectively acts as the GL component liaison and arbiter. 
class Radar extends Element{
    constructor(override_id = ""){
        super("Radar", "radar", override_id)
        this.key ="Panel"
        this.updated = true;
        this.view = "velmax2D"

        //Canvas properties
        this.height = 1000;
        this.width = 1000;
        this.scale = [1.0,1.0];

        //2D Graph Application Properties
        this.x_scale = 0;
        this.y_scale = 0;
        this.x_max = 0;
        this.y_max = 0;
        this.x_min = 0;
        this.y_min = 0;
        this.x_coord_ticks = 0;

        this.canvas = null;
        this.ctx = null;
        this.interval = null;
        this.frame_index = 0;

        this.fft = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.adcA1_I = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.adcA1_Q = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.adcA2_I = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.adcA2_Q = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.adcB_I = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.adcB_Q = new Swapbuffer(new Int16Array(256), new Int16Array(256));
        this.pdat = new Swapbuffer(new Int16Array(12), new Int16Array(12));
        this.ddat = new Swapbuffer(new Int8Array(6), new Int8Array(6));
        this.vel_max = new Swapbuffer(new Int16Array(180), new Int16Array(180));
        this.range_max = new Swapbuffer(new Int16Array(180), new Int16Array(180));
        this.rpst = new Int8Array(42);


        //Airxt Flight Object Properties
        this.vel_club_lateral_prior = 0;
        this.vel_club_lateral_post = 0;
        this.vel_impulse_time = 0;
        this.vel_club_radial_prior = 0;
        this.vel_club_radial_post = 0;
        this.x_rpm = 0;
        this.y_rpm = 0;
        this.z_rpm = 0;
    }


    update(element, props){
        props.forEach(prop => {
            let name = prop.Name;
            if(name === "fft"){
                this.fft.swap();
                this.fft.back = new Int16Array(prop.value);
            }else if(name === "adc"){
                this.adc.swap();
                this.adc.back = new Int16Array(prop.value);
            }else if(name === "pdat"){
                this.pdat.swap();
                this.pdat.back = new Int16Array(prop.value);
            }else if(name === "ddat"){
                this.ddat.swap();
                this.ddat.back = new Int8Array(prop.value);
            }else if(name === "rpst"){
                this.rpst = new Int8Array(prop.value);
            }else if(name == "view"){
                this.view = prop.value;
            }else if(name == "height"){
                this.height = prop.value;
            }else if(name == "width"){
                this.width = prop.value;
            }
        });
    }

    mount(){
        //initiate opengl canvas and context
        super.mount();
        this.canvas = document.getElementById("render-canvas");
        this.ctx = this.canvas.getContext("2d");
        this.image = ctx.createImageData(this.width, this.height);
        this.interval = setInterval(this.draw, 100)
    }


     colorMap(value) {
        // Map the value to a range from 0-1
        let t = value / 100;
      
        // Define the heatmap colors
        let colors = [
          [0, 0, 255],   // blue
          [0, 255, 0],   // green
          [255, 255, 0], // yellow
          [255, 0, 0]    // red
        ];
      
        // Find the two closest colors to the mapped value
        let index1 = Math.floor(t * (colors.length - 1));
        let index2 = index1 + 1;
        let weight2 = (t - index1 / (colors.length - 1)) * (colors.length - 1);
        let weight1 = 1 - weight2;
      
        // Interpolate between the two closest colors
        let color = [
          Math.floor(colors[index1][0] * weight1 + colors[index2][0] * weight2),
          Math.floor(colors[index1][1] * weight1 + colors[index2][1] * weight2),
          Math.floor(colors[index1][2] * weight1 + colors[index2][2] * weight2),
          255
        ];

        let heatmap = new Uint8ClampedArray(color);
        return heatmap;
      }

    draw(){ 
        if(this.view === "velmax2D"){
            //velmax2D plots a scalar value as a color map over the field (x,y) where is the time domain
            //and y is the angle domain, we can view this as a scan line approach without clearing the canvas
            let ppi = 180/this.height; //pixels per index
            let i = 0.0;
            let image = ctx.createImageData(1, this.height); //draw 180 degree elements across the amount of pixels
            this.vel_max.front().forEach((value, index) => {   
                let pixel_index = Math.floor(i)
                let color = this.colorMap(value)
                 //uint8clampedarray data
                 image.data[pixel_index*this.canvas.width*4] = color[0];
                 image.data[pixel_index*this.canvas.width*4+ 1] = color[1];
                 image.data[pixel_index*this.canvas.width*4 + 2] = color[2];
                 image.data[pixel_index*this.canvas.width*4 + 3] = color[3];
                i += ppi;
            })
            this.ctx.putImageData(line, frame_index++, 0);

            if(frame_index > this.width){
                frame_index = 0;
                this.ctx.clearRect(0,0,this.canvas.width, this.canvas.height);
            }
        }else if(this.view === "reflmax2D"){
            //reflmax2D plots a scalar value as a color map over the field (x,y) where is the time domain
            //and y is the angle domain, we can view this as a scan line approach without clearing the canvas
            let ppi = 180/this.height; //pixels per index
            let i = 0.0;
            let image = ctx.createImageData(1, this.height); //draw 180 degree elements across the amount of pixels
            this.range_max.front().forEach((value, index) => {   
                let pixel_index = Math.floor(i)
                let color = this.colorMap(value)
                //uint8clampedarray data
                image.data[pixel_index*this.canvas.width*4] = color[0];
                image.data[pixel_index*this.canvas.width*4+ 1] = color[1];
                image.data[pixel_index*this.canvas.width*4 + 2] = color[2];
                image.data[pixel_index*this.canvas.width*4 + 3] = color[3];
                i += ppi;
            })
            this.ctx.putImageData(line, frame_index++, 0);

            if(frame_index > this.width){
                frame_index = 0;
                this.ctx.clearRect(0,0,this.canvas.width, this.canvas.height);
            }
        }else if(this.view === "fftraw"){
            //fftraw plots a scalar value as a color map over the field (x,y) where is the time domain
            //and y is the angle domain, we can view this as a scan line approach without clearing the canvas
            let image = ctx.createImageData(this.width, this.height);
            this.fft.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })
            this.ctx.putImageData(image, 0, 0);
        }else if(this.view === "adcraw"){
            //adcraw plots a scalar value as a color map over the field (x,y) where is the time domain
            //and y is the angle domain, we can view this as a scan line approach without clearing the canvas
            let image0 = ctx.createImageData(this.width, this.height);

            this.adcA1_I.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })
            this.adcA1_Q.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })
            this.adcA2_I.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })
            this.adcA2_Q.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })
            this.adcB_I.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })
            this.adcB_Q.front().forEach((value, index) => {    
                let y_pixel_index = Math.floor((value/65535)*this.canvas.height);
                let x_pixel_index = Math.floor((index/256)*this.canvas.width);
                let color = this.colorMap(value/65535);

                //uint8clampedarray data
                image.data[y_pixel_index*this.canvas.width*4 + x_pixel_index*4] = color[0];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 1] = color[1];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 2] = color[2];
                image.data[y_pixel_index*this.canvas.width + x_pixel_index*4 + 3] = color[3];
            })

            this.ctx.putImageData(image, 0, 0);


        }else if(this.view === "pdatmax2D"){
            //pdatmax2D plots range/azimuth to target given by distance angle. We then indicate target velocity
            //with the target color/size. 
            let index = 0;
            let step = 4;
            let image = ctx.createImageData(this.width, this.height);
            for(index = 0; i < this.pdat.front().length; i+=step){
                let data = this.pdat.front().slice(index, index+step);
                let range = data[0];
                let speed = data[1];
                let azimuth = data[2];
                let db = data[3];

                //here we plot a pixelated blob at the range/azimuth heatmaped by speed and spread by a guassian
                let color = this.colorMap(speed/100);
                let radius = 6.5;
                let max_range = 30.0;
                let pp_range = max_range/this.height;

                // 
            }
        }

    }

    render(){
        //pass updated data to opengl context
        super.render();
        
    }

    html(){
        return `<div id=${this.id} class='canvas'></div>`
    }    
}


class Swapbuffer{
    constructor(present, back){
        this.present = present;
        this.back = back;
    }

    swap(){
        let temp = this.present;
        this.present = this.back;
        this.back = temp;
    }

    front(){
        return this.present;
    }

    back(){
        return this.back;
    }
}
