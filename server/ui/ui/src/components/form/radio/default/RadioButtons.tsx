// vendor
import React, {
    FC,
    ChangeEvent,
    useCallback,
    useRef,
  } from 'react';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faDotCircle, faCircle } from '@fortawesome/free-regular-svg-icons';
// css
import './RadioButtons.scss';


interface RadioValues {
  text: string;
  value: string;
}

interface Props {
  name: string;
  radioType: string;
  updateValue: (arg0: string) => void;
  value: string;
  values: Array<RadioValues>;
  disabledValues?: Array<any>;
}

const RadioButtons: FC<Props> = ({
  name,
  radioType,
  updateValue,
  value,
  values,
  disabledValues,
}:Props) => {
  //refs
  const inputRef = useRef(null);


  /**
  * Method updates the search value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setSearch}
  */
  const handleChangeEvent = useCallback((event: ChangeEvent) => {
    const element = event.target as HTMLInputElement;
    updateValue(element.getAttribute('value'));
    return;
  }, [updateValue])

  const root = document.documentElement;
  const primaryHash = root.style.cssText ?
    root.style.cssText.split(';')[0].split(':')[1]
    : '#957299';

  return (
    <section className="RadioButtons">
      <p><b>Select a {radioType}:</b></p>

      <div className="RadioButtons__buttons">
        {values.map((item) => {
          const isDisabled = disabledValues && (disabledValues.indexOf(item.value) > -1);
          const labelCSS = classNames({
            RadioButtons__label: true,
            'RadioButtons__label--disabled': isDisabled,
          });
          return (
            <label
              className={labelCSS}
              htmlFor={item.value}
              key={item.value}
            >
              <FontAwesomeIcon
                color={primaryHash}
                icon={value === item.value ? faDotCircle : faCircle}
                size="lg"
              />
              <span>{item.text}</span>
              <input
                className="hidden"
                name={name}
                id={item.value}
                onChange={(event) => handleChangeEvent(event)}
                value={item.value}
                disabled={isDisabled}
                type="radio"
              />
            </label>
          )}
        )}
      </div>

    </section>
  )
}


export default RadioButtons;
