#include "cal_7.h"
#include "ui_cal_7.h"

cal_7::cal_7(QString uid,QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::cal_7), cli(settings.value("host").toString(), settings.value("port").toString()), uid(uid)
{
    ui->setupUi(this);
    loadEventsForWeek();

}

cal_7::~cal_7()
{
    delete ui;
}

void cal_7::eventDeleted(event_entry* ev)
{
    bool ok = cli.deleteEvent(uid, ev->getData().ID);
    if (!ok) {
        QMessageBox::warning(this, "Ошибка", "При удалении события произошла ошибка");
    } else {
        delete ev;
    }
    loadEventsForWeek();
}

void cal_7::loadEventsForWeek()
{
    QDate today = QDate::currentDate();
    QDate startOfWeek = today.addDays(-today.dayOfWeek() + 1);

    QScrollArea* scrollAreas[] = {
        ui->scrollArea_1, ui->scrollArea_2, ui->scrollArea_3,
        ui->scrollArea_4, ui->scrollArea_5, ui->scrollArea_6,
        ui->scrollArea_7
    };

    QLabel* dateLabels[] = {
        ui->date1, ui->date2, ui->date3,
        ui->date4, ui->date5, ui->date6,
        ui->date7
    };

    for (int i = 0; i < 7; ++i) {
        QWidget* oldWidget = scrollAreas[i]->widget();
        if (oldWidget) {
            delete oldWidget;
        }
    }

    for (int i = 0; i < 7; ++i) {
        QDate currentDay = startOfWeek.addDays(i);

        dateLabels[i]->setText(currentDay.toString("dd.MM"));
        if (currentDay == today) {
            dateLabels[i]->setStyleSheet("font-weight: bold; color: blue;");
        } else {
            dateLabels[i]->setStyleSheet("");
        }
        QVector<EventData> events = cli.getEventsByDay(currentDay, uid);
        QWidget* scrollContent = new QWidget;
        QVBoxLayout* layout = new QVBoxLayout(scrollContent);

        if (events.isEmpty()) {
            QLabel* noEventsLabel = new QLabel("Нет событий");
            noEventsLabel->setAlignment(Qt::AlignCenter);
            layout->addWidget(noEventsLabel);
        } else {
            for (const EventData& event : events) {
                event_entry* entry = new event_entry(event, uid);
                connect(entry, &event_entry::deleted, this, &cal_7::eventDeleted);
                connect(entry, &event_entry::edited, this, &cal_7::eventEdited);
                layout->addWidget(entry);
            }
        }
        scrollContent->setLayout(layout);
        scrollAreas[i]->setWidget(scrollContent);
    }
}
void cal_7::showEvent(QShowEvent *event) {
    QWidget::showEvent(event);
    loadEventsForWeek();
}
void cal_7::eventEdited(event_entry* ev)
{

    loadEventsForWeek();
}
