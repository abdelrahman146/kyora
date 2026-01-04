#!/usr/bin/env node

/**
 * Auto-generates field component API reference for forms.instructions.md
 * Extracts prop types from portal-web/src/components/form/ TypeScript files
 *
 * Usage: node scripts/generate-field-docs.js
 * Output: markdown tables with component props, types, defaults, descriptions
 */

import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";
import ts from "typescript";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const FORM_COMPONENTS_DIR = path.join(
  __dirname,
  "../portal-web/src/components/form"
);
const OUTPUT_FILE = path.join(
  __dirname,
  "../.github/instructions/forms-field-reference.md"
);

// Field components to document
const FIELD_COMPONENTS = [
  "TextField.tsx",
  "PasswordField.tsx",
  "SelectField.tsx",
  "CheckboxField.tsx",
  "DateField.tsx",
  "TimeField.tsx",
  "DateRangeField.tsx",
  "FieldArray.tsx",
  "FileUploadField.tsx",
  "ImageUploadField.tsx",
];

/**
 * Extract props interface from TypeScript file
 */
function extractPropsInterface(filePath) {
  const sourceCode = fs.readFileSync(filePath, "utf-8");
  const sourceFile = ts.createSourceFile(
    filePath,
    sourceCode,
    ts.ScriptTarget.Latest,
    true
  );

  let propsInterface = null;
  let propsTypeName = null;

  // Find Props interface (e.g., TextFieldProps, SelectFieldProps)
  function visit(node) {
    if (ts.isInterfaceDeclaration(node)) {
      const name = node.name.text;
      if (name.endsWith("Props")) {
        propsTypeName = name;
        propsInterface = node;
      }
    }

    if (ts.isTypeAliasDeclaration(node)) {
      const name = node.name.text;
      if (name.endsWith("Props")) {
        propsTypeName = name;
        propsInterface = node;
      }
    }

    ts.forEachChild(node, visit);
  }

  visit(sourceFile);

  if (!propsInterface) {
    return null;
  }

  const props = [];

  // Extract props from interface members
  if (ts.isInterfaceDeclaration(propsInterface)) {
    propsInterface.members.forEach((member) => {
      if (ts.isPropertySignature(member)) {
        const propName = member.name?.getText(sourceFile);
        const propType = member.type?.getText(sourceFile);
        const optional = !!member.questionToken;
        const jsDocComment = getJsDocComment(member, sourceFile);

        props.push({
          name: propName,
          type: propType,
          optional,
          description: jsDocComment,
        });
      }
    });
  }

  // Extract props from type alias (if extends other types)
  if (ts.isTypeAliasDeclaration(propsInterface)) {
    const type = propsInterface.type;
    if (ts.isIntersectionTypeNode(type)) {
      type.types.forEach((typeNode) => {
        if (ts.isTypeLiteralNode(typeNode)) {
          typeNode.members.forEach((member) => {
            if (ts.isPropertySignature(member)) {
              const propName = member.name?.getText(sourceFile);
              const propType = member.type?.getText(sourceFile);
              const optional = !!member.questionToken;
              const jsDocComment = getJsDocComment(member, sourceFile);

              props.push({
                name: propName,
                type: propType,
                optional,
                description: jsDocComment,
              });
            }
          });
        }
      });
    }
  }

  return { propsTypeName, props };
}

/**
 * Get JSDoc comment for a node
 */
function getJsDocComment(node, sourceFile) {
  const jsDoc = ts.getJSDocCommentsAndTags(node);
  if (jsDoc.length > 0) {
    const comment = jsDoc[0];
    if (ts.isJSDoc(comment)) {
      return comment.comment?.toString() || "";
    }
  }
  return "";
}

/**
 * Generate markdown table for component props
 */
function generateMarkdownTable(componentName, propsData) {
  if (!propsData) {
    return `### ${componentName}\n\n*No props interface found*\n\n`;
  }

  const { propsTypeName, props } = propsData;

  let markdown = `### ${componentName}\n\n`;
  markdown += `**Interface:** \`${propsTypeName}\`\n\n`;

  if (props.length === 0) {
    markdown += "*No documented props*\n\n";
    return markdown;
  }

  markdown += `| Prop | Type | Required | Description |\n`;
  markdown += `|------|------|----------|-------------|\n`;

  props.forEach((prop) => {
    const required = prop.optional ? "No" : "**Yes**";
    const description = prop.description || "-";
    // Escape pipe characters in type definitions
    const type = prop.type.replace(/\|/g, "\\|");

    markdown += `| \`${prop.name}\` | \`${type}\` | ${required} | ${description} |\n`;
  });

  markdown += "\n";

  return markdown;
}

/**
 * Main function
 */
function main() {
  console.log("🔍 Scanning form components...\n");

  let output = `# Form Field Component API Reference\n\n`;
  output += `*Auto-generated from TypeScript definitions*\n`;
  output += `*Last updated: ${new Date().toISOString()}*\n\n`;
  output += `---\n\n`;
  output += `## Field Components\n\n`;

  FIELD_COMPONENTS.forEach((componentFile) => {
    const componentPath = path.join(FORM_COMPONENTS_DIR, componentFile);

    if (!fs.existsSync(componentPath)) {
      console.log(`⚠️  ${componentFile} not found, skipping`);
      return;
    }

    console.log(`📄 Processing ${componentFile}...`);

    const componentName = componentFile.replace(".tsx", "");
    const propsData = extractPropsInterface(componentPath);
    const markdown = generateMarkdownTable(componentName, propsData);

    output += markdown;
  });

  output += `---\n\n`;
  output += `## Usage Notes\n\n`;
  output += `- All field components require \`name\` prop for form field binding\n`;
  output += `- \`label\` prop is required unless \`hideLabel={true}\` is set\n`;
  output += `- Use within TanStack Form's \`<form.Field>\` wrapper for full validation\n`;
  output += `- See forms.instructions.md for complete integration examples\n\n`;
  output += `## Updating This Reference\n\n`;
  output += `Run: \`node scripts/generate-field-docs.js\`\n\n`;
  output += `To integrate changes into forms.instructions.md, copy relevant sections.\n`;

  // Write output file
  fs.writeFileSync(OUTPUT_FILE, output, "utf-8");

  console.log(`\n✅ Generated: ${OUTPUT_FILE}`);
  console.log(`📊 Total components documented: ${FIELD_COMPONENTS.length}\n`);
}

main();
