import { Component, ViewChild, ElementRef } from '@angular/core';

@Component({
    selector: `card-view`,
    template: `
        <div>
            <md-card>

                <md-card-header>
                    <md-card-title #title>
                        <ng-content select="[title]"></ng-content>
                    </md-card-title>
                </md-card-header>

                <md-card-content #content>
                    <ng-content select="[content]"></ng-content>
                </md-card-content>

            </md-card>
        </div>
    `
})
export class CardComponent {
    @ViewChild('title') titleView: ElementRef
    @ViewChild('content') contentView: ElementRef

    ngAfterViewInit() {
        if (this.titleView.nativeElement.children.length == 0) {
            this.titleView.nativeElement.style.display = 'none'
        }
        if (this.contentView.nativeElement.children.length === 0) {
            this.contentView.nativeElement.style.display = 'none'
        }
    }

}