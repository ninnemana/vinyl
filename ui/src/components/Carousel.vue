<template>
    <div class="carousel">
        <!-- <div class="main">
            <img v-bind:src='mainImage.resource_url'>
        </div> -->
        <v-list class="alternate">
            <div v-bind:key='image.resource_url' v-for='image in images'>
                <a href="javascript:void(0)" @click="setMain">
                    <img v-bind:src='image.resource_url'>
                </a>
            </div>
        </v-list>
    </div>
</template>

<script lang='ts'>
import { Vue, Component, Prop } from 'vue-property-decorator';

@Component
export default class Carousel extends Vue {
    @Prop() readonly images!: any[];
    private mainImage!: any;

    // constructor() {
    //     super();
        
    //     this.images = [];
    // }

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
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.carousel {
    max-width: 600px;
    display: flex;
    align-items: center;
    flex-direction: column;
}

.alternate {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;

    > div {
        width: 60px;
        height: 60px;
        overflow: hidden;
        margin: 10px;
        > img {
            width: 60px;
        }
    }
}
</style>
