<template>
    <div class='pagination-block'>
        <transition-group name="list" tag="div" mode="out-in">
            <a href="javascript:void(0);" v-bind:key='"prev"' @click='change("previous")' class="list-item">
                Prev
            </a>
            <a href="javascript:void(0);" @click='change(n)' v-for="n in listLimit()" v-bind:key="n+1" class="list-item">
                {{ n+1 }}
            </a>
            <a href="javascript:void(0);"  v-bind:key='"next"' @click='change("next")' class="list-item">
                Next
            </a>
        </transition-group>
    </div>
</template>

<script lang='ts'>
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Pagination extends Vue {
    @Prop() private settings: any = {};
    @Prop() private onChange!: (i: number) => void;

    private change(i: number) {
        if (!this.onChange) {
            return;
        }

        this.onChange(i);
    }

    private listLimit() {
        const settings = JSON.parse(JSON.stringify(this.settings)) || {};
        if (!settings.items || !settings.per_page) {
            return [];
        }

        const items = +settings.items;
        const perPage = +settings.per_page;
        const results = new Array(Math.floor(items / perPage));

        for (let index = 0; index < Math.floor(items / perPage); index++) {
            results[index] = index;
        }
        return results;
    }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.pagination-block {
    > div {
        display: flex;
        flex-direction: row;
        margin: 10px;

        > a {
            margin: 10px;
            color: black;
        }
    }
}

.list-enter-active, .list-leave-active {
  transition: all 1s;
}
.list-enter, .list-leave-to /* .list-leave-active below version 2.1.8 */ {
  opacity: 0;
  transform: translateY(30px);
}
</style>
