/* 
* Licensed Materials - Property of IBM.
* Copyright IBM Corporation 2023. All Rights Reserved.
* U.S.Government Users Restricted Rights - Use, duplication or disclosure
* restricted by GSA ADP Schedule Contract with IBM Corp.
*
* Contributors:
* IBM Corporation - initial API and implementation
*/

export class WaziCardEditors {
    static nameMarker = 'wazi';
    static editorIds = ['che-code'];
    /**
     * Updates the specific card's for available editors
     * Removes all options besides VS Code
     * @param cardsArray array of cards returned from buildCardsList()
     */
    public static update(cardsArray: React.ReactElement<any, string | React.JSXElementConstructor<any>>[]): React.ReactElement<any, string | React.JSXElementConstructor<any>>[] {
        cardsArray.forEach(function (card, index, cardsArray) {
            if (card.props.metadata.displayName.toLowerCase().includes(WaziCardEditors.nameMarker)) {
                cardsArray[index].props.targetEditors = card.props.targetEditors.filter(editor => 
                    WaziCardEditors.editorIds.some(editorId => 
                        editor.id.includes(editorId)));
            }
        });
        return cardsArray
    }
}