<template>
    <div class="header">
        <div class="carousel">
            <div class="main">
                <v-img aspect-ratio="2" v-bind:src="mainImage.resource_url" contain></v-img>
            </div>
            <v-row justify="space-around">
                <v-col cols="images.length">
            <!-- <v-list class="alternate"> -->
                <div 
                    v-bind:key='image.resource_url' 
                    v-for='image in images'
                    v-bind:class="{ 
                        active: mainImage.resource_url === image.resource_url
                    }" 
                >
                    <a href="javascript:void(0)" @click="setMain(image)">
                        <v-img src="https://picsum.photos/510/300?random" aspect-ratio="1.7"></v-img>
                        <v-img aspect-ratio="2" v-bind:src="image.resource_url" contain></v-img>
                    </a>
                </div>
                </v-col>
            </v-row>
            <!-- </v-list> -->
        </div>
        <div class="info">
            <div class="artists">
                
            </div>
        </div>
    </div>
</template>

<script lang='ts'>
import { Vue, Component, Prop } from 'vue-property-decorator';

@Component
export default class Carousel extends Vue {
    @Prop() readonly images!: any[];
    private mainImage!: any;

    constructor() {
        super();
        
        const primary = this.images.filter((v: any) => v.type === 'primary');
        if (primary.length > 0) {
            this.mainImage = primary[0];
        }
    }

    computed() {
        const def = {
            images: this.images,
            mainImage: {},
        };

        const primary = this.images.filter((v: any) => v.type === 'primary');
        if (primary.length > 0) {
            def.mainImage = primary[0];
        }

        return def;
    }

    setMain(image: any) {
        if (image) {
            this.mainImage = image;
        }
    }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.carousel {
    max-width: 600px;
    display: flex;
    align-items: center;
    flex-direction: column;

    .main {
        width: 400px;
        height: 400px;
        margin-bottom: 20px;
        overflow: hidden;

        img {
            max-width: 400px;
            object-fit: cover;
        }
    }
}

// .alternate {
//     display: flex;
//     flex-direction: row;
//     flex-wrap: wrap;

//     > div {
//         margin: 10px;
//         position: relative;
//         padding: 4px;
//         border: 2px solid transparent;
//         display: flex;

//         > a {
//             height: 60px;
//             overflow: hidden;

//             > img {
//                 width: 60px;
//                 object-fit: cover;
//             }
//         }
//     }

//     > div.active {
//         opacity: 80%;
//         z-index: 10;
//         border-radius: 0.3rem;
//         border-color: #373131;
//     }
// }
</style>
