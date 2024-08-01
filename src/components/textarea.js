import "./style.css";

function TextAreaBox(props) {
  return (
    <div className="input-box" style={{ marginTop: props.margin }}>
      <label>{props.label}</label>
      <br />
      <textarea 
      style={{
        maxLines:props.max,
        minHeight:props.minHeight
      }}  
      
      type={props.type} placeholder={props.placeholder} />
    </div>
  );
}

export default TextAreaBox;
