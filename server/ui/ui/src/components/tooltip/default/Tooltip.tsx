// vendor
import React,
{
  FC,
  useCallback,
  useEffect,
  useRef,
  useState,
  MouseEvent
} from 'react';
import classNames from 'classnames';
import ReactTooltip from 'react-tooltip';
// assets
import './Tooltip.scss';

interface Props {
  isVisible: boolean;
  text: string;
  updateIsVisible?: () => void;
}

const Tooltip: FC<Props> = ({
  isVisible,
  text,
  updateIsVisible,
}: Props) => {
  const tooltipRef = useRef(null);
  const [tooltipVisibility, setTooltipVisibility] = useState(false);
  const [tooltipExpanded, updateExpanded] = useState(false);

  if (tooltipVisibility !== isVisible) {
    setTooltipVisibility(isVisible);
  }

  /**
    *  @param {String} section
    *  @param {Boolean} tooltipExpanded
    *  @param {Function} updateExpanded
    *  closes tooltip box when tooltip is open and the tooltip has not been clicked on
    *
  */
  const hideTooltip = (event: Event) => {
    if (tooltipExpanded
        && !tooltipRef.current.contains(event.target)
    ) {
      updateExpanded(false);
    }
    if (updateIsVisible) {
      updateIsVisible();
    }
    ReactTooltip.hide(tooltipRef.current);
  }


  /**
    *  @param {Event} event
    *  shows tooltip box when clicked
  */
  const showToolTip = useCallback(() => {
    ReactTooltip.show(tooltipRef.current);
  }, [])

  useEffect(() => {
    ReactTooltip.show(tooltipRef.current);

    window.addEventListener('click', hideTooltip);
    return () => {
      window.removeEventListener('mouseUp', hideTooltip);
    }
  }, []);



  useEffect(() => {
    if (isVisible && text) {
      showToolTip();
    }
  }, [text, isVisible, showToolTip]);

  // declare css here
  const toolTipCSS = classNames({
    Tooltip: true,
    hidden: !tooltipVisibility,
  });

  return (
    <div className={toolTipCSS}>
      <div
        data-event="click focus"
        data-tip={text}
        ref={tooltipRef}
        role="presentation"
      />
      <ReactTooltip
        place="bottom"
        offset={{left: 20, bottom: 5}}
      />

    </div>
  );
}

export default Tooltip;
