// vendor
import React, { FC, MouseEvent, useRef } from 'react';
import classNames from 'classnames';
// hooks
import { useEventListener } from 'Hooks/events/index';
// css
import './Dropdown.scss';


interface ListObject {
  name: string;
  type?: string;
}

interface Payload {
  selectedItem?: ListObject;
}

interface Action {
  type: string;
  payload?: Payload;
}

interface Props {
  customStyle: string,
  dispatch: (arg: Action) => void;
  label: string,
  listItems: Array<ListObject>,
  visibility: boolean,
}

const Dropdown: FC<Props> = ({
  customStyle,
  dispatch,
  label,
  listItems,
  visibility,
}: Props) => {

    const dropdownRef = useRef(null);

    /**
    * Method handles click on a individual item in the dropdown menu and dispatches set to the parent
    * @param {MouseEven} event
    * @param {string} item
    * @return {void}
    * @calls {event#stopPropagation}
    * @calls {dispatch}
    */
    const clickItem = (event: MouseEvent, item: ListObject) => {
      event.stopPropagation();

      dispatch({
        type: 'set',
        payload: { selectedItem: item }
      });
    }


    /**
    * Method provides a way for child componts to update state
    * @param {Object} evt
    * @fires updateCodeFileUrl
    */
    const globalClickHandler =  (event: Event) => {
        if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
          dispatch({ type: 'close' });
        }
    };

    // Add event listener using our hook
    useEventListener('click', globalClickHandler);

    /**
    * Method handles dropdown click to open the dropdown menu
    * @param {MouseEven} event
    * @return {void}
    * @calls {event#stopPropagation}
    * @calls {dispatch}
    */
    const clickDropdown = (event: MouseEvent) => {
      event.stopPropagation();
      dispatch({ type: 'toggle' });
    }

    const dropdownCSS = classNames({
      'Dropdown relative': true,
      'Dropdown--open': visibility,
      'Dropdown--collapsed': !visibility,
      [`Dropdown--${customStyle}`]: customStyle,
    });

    return (
      <div
        className="relative"
        ref={dropdownRef}
        role="presentation"
      >
        <button
          className={dropdownCSS}
          onClick={clickDropdown}
        >
        {label}
        </button>
        {
          visibility && (
            <menu className="Dropdown__menu">
              {
                listItems.map(item => (
                  <button
                    className="Dropdown__item"
                    key={item.name}
                    onClick={(evt) => { clickItem(evt, item); }}
                    role="presentation"
                  >
                    <span>{item.name}</span>
                    <span>{item.type ? ` (${item.type})`: ''}</span>

                  </button>
                ))
              }
            </menu>
          )
        }
      </div>
    );
}


export default Dropdown;
