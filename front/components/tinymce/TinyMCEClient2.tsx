import React from "react";
import { Editor } from "@tinymce/tinymce-react";

interface Props {
  value: string;
}

const TinyEditor2 = ({ value }: Props) => {
  return (
    <Editor
      apiKey={process.env.NEXT_PUBLIC_TINYMCE_API_KEY}
      initialValue={value}
      init={{
        height: 360,
        plugins: [
          // Core editing features
          "anchor",
          "autolink",
          "charmap",
          "codesample",
          "emoticons",
          "link",
          "lists",
          "media",
          "searchreplace",
          "table",
          "visualblocks",
          "wordcount",
          // Your account includes a free trial of TinyMCE premium features
          // Try the most popular premium features until Aug 28, 2025:
          "checklist",
          "mediaembed",
          "casechange",
          "formatpainter",
          "pageembed",
          "a11ychecker",
          "tinymcespellchecker",
          "permanentpen",
          "powerpaste",
          "advtable",
          "advcode",
          "advtemplate",
          "ai",
          "uploadcare",
          "mentions",
          "tinycomments",
          "tableofcontents",
          "footnotes",
          "mergetags",
          "autocorrect",
          "typography",
          "inlinecss",
          "markdown",
          "importword",
          "exportword",
          "exportpdf",
        ],
        toolbar:
          "undo redo | blocks fontfamily fontsize | bold italic underline strikethrough | link media table mergetags | addcomment showcomments | spellcheckdialog a11ycheck typography uploadcare | align lineheight | checklist numlist bullist indent outdent | emoticons charmap | removeformat",
        tinycomments_mode: "embedded",
        tinycomments_author: "Author name",
        mergetags_list: [
          { value: "First.Name", title: "First Name" },
          { value: "Email", title: "Email" },
        ],
        ai_request: (
          _request: unknown,
          respondWith: { string: (fn: () => Promise<string>) => void }
        ) =>
          respondWith.string(() =>
            Promise.reject("See docs to implement AI Assistant")
          ),
        uploadcare_public_key: "2f564a356ca8519fe277",
      }}
    />
  );
};

export default TinyEditor2;
