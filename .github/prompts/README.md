# Kyora Prompt Files

Reusable prompt files for common development tasks in the Kyora monorepo.

## Available Prompts

### Feature Development

- **`/add-backend-feature`** - Add new feature to backend (Go API)
- **`/add-portal-feature`** - Add new feature to portal-web (React dashboard)

### Project Scaffolding

- **`/create-backend-project`** - Create new backend domain module
- **`/create-frontend-project`** - Create new frontend feature module

### Debugging

- **`/debug-backend`** - Debug and fix backend issues
- **`/debug-portal-web`** - Debug and fix portal-web issues
- **`/fix-ui-ux`** - Fix UI/UX issues in portal-web

## Usage

### In Chat View

Type `/` followed by the prompt name in the chat input:

```
/add-backend-feature Add endpoint to export orders as CSV
```

```
/fix-ui-ux Button alignment broken in RTL mode on mobile
```

### From Command Palette

1. Open Command Palette (`⌘⇧P` / `Ctrl+Shift+P`)
2. Run `Chat: Run Prompt`
3. Select prompt from list
4. Enter required parameters

### From Prompt File

1. Open prompt file in editor
2. Click play button in title bar
3. Choose to run in current or new chat session

## Prompt Structure

Each prompt file includes:

- **YAML Frontmatter**: Description, agent, tools, model configuration
- **Feature Requirements**: Input variable for user to describe their task
- **Instructions**: Links to relevant instruction files
- **Standards**: Quality and architecture requirements
- **Workflow**: Step-by-step implementation guide
- **Done Criteria**: Completion checklist

## Tips

1. **Be Specific**: Provide detailed feature/issue descriptions in the input
2. **Reference Examples**: Prompts suggest similar existing features to reference
3. **Follow Standards**: Prompts enforce Kyora architecture patterns
4. **Test Thoroughly**: Prompts include testing steps in workflow
5. **Use Instruction Files**: Prompts link to specialized instruction files for details

## Instruction Files Reference

Prompts reference these instruction files (SSOT):

- `backend-core.instructions.md` - Backend architecture
- `backend-testing.instructions.md` - Backend testing
- `portal-web-architecture.instructions.md` - Frontend architecture
- `portal-web-development.instructions.md` - Frontend workflow
- `forms.instructions.md` - Form system
- `ui-implementation.instructions.md` - UI patterns
- `design-tokens.instructions.md` - Design tokens
- `charts.instructions.md` - Data visualization
- `ky.instructions.md` - HTTP client
- `stripe.instructions.md` - Billing
- `resend.instructions.md` - Email
- `asset_upload.instructions.md` - File uploads

## Adding New Prompts

1. Create `.prompt.md` file in this directory
2. Add YAML frontmatter with description, agent, tools
3. Use `${input:variableName:placeholder}` for user inputs
4. Reference instruction files with relative links
5. Include workflow steps and done criteria
6. Update this README with new prompt
7. Test prompt with various inputs

## Examples

### Add Backend Feature

```
/add-backend-feature Add webhook endpoint for Stripe payment success events
```

### Fix UI Issue

```
/fix-ui-ux Form fields overlapping on mobile screens in Arabic mode
```

### Debug Backend

```
/debug-backend Getting 500 error when creating order with multiple line items
```

### Create New Module

```
/create-backend-project Create subscription domain for managing recurring billing
```
